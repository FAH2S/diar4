package sharedapi
import (
    "testing"
    "net/http"
    "net/http/httptest"
    "strings"
    "reflect"
)

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
            name:           "wringType",
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


