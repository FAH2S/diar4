package shareddb
import (
    "os"
    "fmt"
    "testing"
    "strings"
    "errors"
    "database/sql"
    "reflect"
)
import (
    "github.com/lib/pq"
)

//{{{ buildConnStrFromEnvFn
// succ
func Test_BuildConnStrFromEnv_Succ(t *testing.T){
    os.Clearenv()
    // Set required env vars for test
    os.Setenv("DB_USER", "testuser")
    os.Setenv("DB_PWD", "testpass")
    os.Setenv("DB_NAME", "testdb")
    os.Setenv("DB_HOST", "localhost")
    os.Setenv("DB_PORT", "5432")

    // expected outcome
    expected := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
    // check
    actual, err := buildConnStrFromEnvFn()
    if err != nil {
        t.Fatalf("Fatal, expected no error, got: %v", err)
    }
    if actual != expected {
        t.Errorf("Wrong return value:\nExpected:\t%s\nGot:\t\t%s", expected, actual)
    }
}


// fail, missing env
func Test_BuildConnStrFromEnv_MissingEnv(t *testing.T){
    // Clear env to cause err
    os.Clearenv()

    // expected outcome
    expected := ""
    expectedErr := "buildConnStrFromEnvFn: Required environment variable DB_USER, not set\n"
    // check
    actual, err := buildConnStrFromEnvFn()
    if err == nil {
        t.Fatalf("Fatal, epected error, got nil")
    }
    if actual != expected {
        t.Errorf("Wrong string value:\nExpected:\t%s\nGot:\t\t%s", expected, actual)
    }
    if err.Error() != expectedErr {
        t.Errorf("Wrong error value:\nExpected:\t%s\nGot:\t\t%s", expectedErr, err.Error())
    }
}
//}}} buildConnStrFromEnvFn


//{{{ HandlePgError
func Test_HandlePgErrorFn(t *testing.T) {
    tests := []struct {
        name                string
        inputErr            error
        expectedStatus      int
        expectedErrSubStr   string
    }{
        {
            name:               "UniqueViolation",
            inputErr:           &pq.Error{Code: "23505"},
            expectedStatus:     409,
            expectedErrSubStr:  "users already exists",
        },{
            name:               "CheckConstraintViolation",
            inputErr:           &pq.Error{Code: "23514"},
            expectedStatus:     422,
            expectedErrSubStr:  "invalid users data/format",
        },{
            name:               "UnknownColumn",
            inputErr:           &pq.Error{Code: "42703"},
            expectedStatus:     400,
            expectedErrSubStr:  "unknown column used",
        },{
            name:               "UnhandledPqError",
            inputErr:           &pq.Error{Code: "999999999"},
            expectedStatus:     500,
            expectedErrSubStr:  "failed to execute query",
        },{
            name:               "NonPqError",
            inputErr:           fmt.Errorf("some non pq error"),
            expectedStatus:     500,
            expectedErrSubStr:  "unexpected error",
        },

    }

    // Iterate
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            statusCode, err := HandlePgErrorFn("users", tc.inputErr)
            if statusCode != tc.expectedStatus {
                t.Errorf("Wrong status code:\nExpected:\t%d\nGot:\t\t%d", tc.expectedStatus, statusCode)
            }
            if err == nil || !strings.Contains(err.Error(), tc.expectedErrSubStr) {
                t.Errorf("Wrong error:\nExpected:\t%q\nGot:\t\t%q", tc.expectedErrSubStr, err)
            }
        })
    }
}
//}}} HandlePgError


//{{{ HandleSelectErrorFn
func Test_HandleSelectErrorFn(t *testing.T) {
    tests := []struct{
        name                string
        inputErr            error
        expectedStatusCode  int
        expectedErrSubStr   string
    }{
        {
            name:               "success",
            inputErr:           nil,
            expectedStatusCode: 200,
            expectedErrSubStr:  "",
        }, {
            name:               "notFound",
            inputErr:           sql.ErrNoRows,
            expectedStatusCode: 404,
            expectedErrSubStr:  "HandleSelectErrorFn: user not found/dosen't exist",
        }, {
            name:               "queryNotExecuted",
            inputErr:           fmt.Errorf("any error"),
            expectedStatusCode: 500,
            expectedErrSubStr:  "failed to execute query",
        },
    }
    // Iterate
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T){
            statusCode, err := HandleSelectErrorFn(tc.inputErr)
            // Check status code
            if statusCode != tc.expectedStatusCode {
                t.Errorf("Wrong status code:\nExpected:\t%d\nGot:\t\t%d", tc.expectedStatusCode, statusCode)
            }
            // Check error, not expecting but got err
            if tc.expectedErrSubStr == "" && err != nil {
                t.Errorf("Wrong error:\nExpected:\t%q\nGot:\t\t%q", tc.expectedErrSubStr, err)
            }
            // Check error, expecting but got none
            if tc.expectedErrSubStr != "" && err == nil {
                t.Errorf("Wrong error:\nExpected:\t%q\nGot:\t\t%v", tc.expectedErrSubStr, err)
            }
            // Chech error, substring match
            if tc.expectedErrSubStr != "" && err != nil && !strings.Contains(err.Error(), tc.expectedErrSubStr) {
                t.Errorf("Wrong error:\nExpected:\t%q\nGot:\t\t%q", tc.expectedErrSubStr, err)
            }
        })
    }
}
//}}} HandleSelectErrorFn


//{{{ CheckRows


//{{{ helper
type mockResult struct {
    rows    int64
    err     error
}
func (m mockResult) RowsAffected() (int64, error) {
    return m.rows, m.err
}
// Not used but still need to be implemented /shrug
func (m mockResult) LastInsertId() (int64, error) {
    return 0, nil
}
//}}} helper


//{{{ CheckRowsAffectedInsert
func Test_CheckRowsAffectedInsertFn(t *testing.T) {
    tests := []struct {
        name        string
        result      sql.Result
        expectError bool
        errorSubstr string
    }{
        {
            name:        "Success",
            result:      mockResult{rows: 1, err: nil},
            expectError: false,
        }, {
            name:        "WrongRowCount",
            result:      mockResult{rows: 0, err: nil},
            expectError: true,
            errorSubstr: "expected 1 row affected",
        }, {
            name:        "ErrorFromRowsAffected",
            result:      mockResult{rows: 0, err: errors.New("boom")},
            expectError: true,
            errorSubstr: "failed to check rows affected",
        },
    }
    // Iterate
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            err := CheckRowsAffectedInsertFn(tc.result)
            if tc.expectError {
                if err == nil || !strings.Contains(err.Error(), tc.errorSubstr) {
                    t.Errorf("Wrong error:\nExpected:\t%q\nGot:\t\t%v", tc.errorSubstr, err)
                }
            } else {
                if err != nil {
                    t.Errorf("Wrong error:\nExpected:\t%q\nGot:\t\t%v", tc.errorSubstr, err)
                }
            }
        })
    }
}
//}}} CheckRowsAffectedInsert


//{{{ CheckRowsAffectedFn
func Test_CheckRowsAffectedFn(t *testing.T) {
    tests := []struct {
        name                string
        result              sql.Result
        expectedErrSubStr   string
        expectedStatusCode  int
    }{
        {
            name:               "Success",
            result:             mockResult{rows: 1, err: nil},
            expectedErrSubStr:  "",
            expectedStatusCode: 200,
        }, {
            name:               "FailNotFound",
            result:             mockResult{rows: 0, err: nil},
            expectedErrSubStr:  "no rows were affected",
            expectedStatusCode: 404,
        }, {
            name:               "FailToCheckRows",
            result:             mockResult{rows: 0, err: errors.New("failed to check rows")},
            expectedErrSubStr:  "failed to check rows affected:",
            expectedStatusCode: 500,
        },
    }
    // Iterate
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            statusCode, err := CheckRowsAffectedFn(tc.result)
            // Check status code
            if tc.expectedStatusCode != statusCode {
                t.Errorf("Wrong statusCode:\nExpected:\t%d\nGot:\t\t%d", tc.expectedStatusCode, statusCode)
            }
            // Chech error, substring match
            if tc.expectedErrSubStr != "" && err != nil && !strings.Contains(err.Error(), tc.expectedErrSubStr) {
                t.Errorf("Wrong error:\nExpected:\t%q\nGot:\t\t%q", tc.expectedErrSubStr, err)
            }
        })
    }
}

//}}} CheckRowsAffectedFn


//}}} CheckRows


//{{{ BuildSetPartsFn
func Test_BuildSetPartsFn(t *testing.T) {
    tests := []struct {
        name                string
        inputData           map[string]interface{}
        expectedSetParts    []string
        expectedArgs        []interface{}
        expectedErrSubStr   string
    }{
        {
            name:               "successUpdateUserHash",
            inputData:          map[string]interface{}{
                "hash":"hash_string",
            },
            expectedSetParts:   []string{"hash = $1"},
            expectedArgs:       []interface{}{"hash_string"},
            expectedErrSubStr:  "",
        }, {
            name:               "successUpdateUserHashSalt",
            inputData:          map[string]interface{}{
                "hash":"hash_string",
                "salt":"salt_string",
            },
            expectedSetParts:   []string{"hash = $1", "salt = $2"},
            expectedArgs:       []interface{}{"hash_string", "salt_string"},
            expectedErrSubStr:  "",
        }, {
            name:               "failUpdateUser",
            inputData:          map[string]interface{}{
            },
            expectedSetParts:   nil,
            expectedArgs:       nil,
            expectedErrSubStr:  "no fields to update",
        },
    }
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            setParts, args, err := BuildSetPartsFn(tc.inputData)
            // Chck setParts
            if !reflect.DeepEqual(tc.expectedSetParts, setParts) {
                t.Errorf("Wrong setParts:\nExpected:\t%v\nGot:\t\t%v", tc.expectedSetParts, setParts)
            }

            // Check args
            if !reflect.DeepEqual(tc.expectedArgs, args) {
                t.Errorf("Wrong args:\nExpected:\t%v\nGot:\t\t%v", tc.expectedArgs, args)
            }

            // Chech error, substring match
            if tc.expectedErrSubStr != "" && err != nil && !strings.Contains(err.Error(), tc.expectedErrSubStr) {
                t.Errorf("Wrong error:\nExpected:\t%q\nGot:\t\t%q", tc.expectedErrSubStr, err)
            }
        })
    }
}

//}}} BuildSetParts


