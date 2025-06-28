package sharedapi
import (
    "testing"
    "net/http"
    "net/http/httptest"
    "strings"
    "reflect"
)

//{{{ Test MapStatusCodeFn
func Test_MapStatusCodeFn(t *testing.T) {
    tests := []struct {
        name            string
        inputStatusCode int
        inputAction     string
        inputEntity     string
        inputName       string
        inputError      error
        ExpectedMsg     string
        ExpectedErrMsg  string
        ExpectedBool    bool
    }{
        {
            name:               "200Success",
            inputStatusCode:    200,
            inputAction:        "create",
            inputEntity:        "user",
            inputName:          "test_user",
            inputError:         nil,
            ExpectedMsg:        "Success: create user 'test_user'",
            ExpectedErrMsg:     "",
            ExpectedBool:       true,
        }, {
            name:               "404Fail",
            inputStatusCode:    404,
            inputAction:        "read",
            inputEntity:        "user",
            inputName:          "test_user",
            inputError:         nil,
            ExpectedMsg:        "Fail: read user 'test_user'",
            ExpectedErrMsg:     "User not found, dosen't exist",
            ExpectedBool:       false,
        }, {
            name:               "500Unknown",
            inputStatusCode:    999,
            inputAction:        "create",
            inputEntity:        "user",
            inputName:          "",
            inputError:         nil,
            ExpectedMsg:        "Fail: create user",
            ExpectedErrMsg:     "Unknown error",
            ExpectedBool:       false,
        }, {
            name:               "500InternalServerError",
            inputStatusCode:    500,
            inputAction:        "create",
            inputEntity:        "user",
            inputName:          "",
            inputError:         nil,
            ExpectedMsg:        "Fail: create user",
            ExpectedErrMsg:     "Internal server error",
            ExpectedBool:       false,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            message, errMessage, success := MapStatusCodeFn(
                tc.inputStatusCode,
                tc.inputAction,
                tc.inputEntity,
                tc.inputName,
                tc.inputError,
            )
            if tc.ExpectedMsg != message{
                t.Fatalf("\nExpected:\t%q\nGot:\t\t%q", tc.ExpectedMsg, message)
            }
            if tc.ExpectedErrMsg != errMessage{
                t.Fatalf("\nExpected:\t%q\nGot:\t\t%q", tc.ExpectedErrMsg, errMessage)
            }
            if tc.ExpectedBool != success{
                t.Fatalf("\nExpected:\t%v\nGot:\t\t%v", tc.ExpectedBool, success)
            }
        })
    }
}
//}}} Test MapStatusCodeFn


//{{{ Test ExtractJSONValueFn
// test: succ, invalid JSON, missing key, target wrong type
func Test_ExtractJSONValueFn(t *testing.T) {
    tests := []struct {
        name            string
        jsonBody        string
        key             string
        target          interface{}
        expectedErr     string
        expectedTarget  interface{}
    }{
        {
            name:           "success",
            jsonBody:       `{"username": "test_user"}`,
            key:            "username",
            target:         new(string),
            expectedErr:    "",
            expectedTarget: "test_user",
        }, {
            name:           "invalidJSON",
            jsonBody:       `{"username": "`,
            key:            "username",
            target:         new(string),
            expectedErr:    "Invalid JSON",
            expectedTarget: "test_user",
        }, {
            name:           "missingKey",
            jsonBody:       `{"username": "test_user"}`,
            key:            "missing_key",
            target:         new(string),
            expectedErr:    "missing key",
            expectedTarget: "test_user",
        }, {
            name:           "wrongType",
            jsonBody:       `{"username": "test_user"}`,
            key:            "username",
            target:         new(int),
            expectedErr:    "failed to unmarshal value",
            expectedTarget: "test_user",
        },

    }
    // iterate
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.jsonBody))
            req.Header.Set("Content-Type", "application/json")
            err := ExtractJSONValueFn(req, tc.key, tc.target)

            // Don't expect error but got one
            if tc.expectedErr == "" && err != nil {
                t.Fatalf("\nExpected:\t%q\nGot:\t\t%v", tc.expectedErr, err)
            }
            // Expect error but either didn't get it or err msg are not same/contain
            if tc.expectedErr != "" && (err == nil || !strings.Contains(err.Error(), tc.expectedErr)) {
                t.Fatalf("\nExpected:\t%q\nGot:\t\t%v", tc.expectedErr, err)
            }
            // Compare targets only if no errors
            if err == nil {
                // makes `got` type interface but with value [string]
                got := reflect.ValueOf(tc.target).Elem().Interface()
                if !reflect.DeepEqual(got, tc.expectedTarget) {
                    t.Fatalf("\nExpected:\t%v\nGot:\t\t%v", tc.expectedTarget, got)
                }
            }
        })
    }

}
//}}} Test ExtractJSONValueFn

