package sharedapi
import (
    "testing"
    "net/http/httptest"
    "encoding/json"
    "reflect"
)

// TODO: down the line more cases need to be added
func Test_WriteJSONResponseFn(t *testing.T) {
    tests := []struct {
        name        string
        statusCode  int
        message     string
        errMsg      string
        data        interface{}
    }{
        {
            name:       "succ",
            statusCode: 201,
            message:    "Success: create user 'test_user'",
            errMsg:     "",
            data:       nil,
        }, {
            name:       "succ",
            statusCode: 400,
            message:    "Fail: to process X",
            errMsg:     "Invalid JSON",
            data:       nil,
        }, {
            name:       "succ",
            statusCode: 200,
            message:    "Success: read user 'test_user'",
            errMsg:     "",
            data:       map[string]interface{}{"foo": "bar"},
        },

    }

    // Iterate
    for _, tc := range tests{
        t.Run(tc.name, func(t *testing.T) {
            w := httptest.NewRecorder()
            WriteJSONResponseFn(w, tc.statusCode, tc.message, tc.errMsg, tc.data)

            // Check status code
            if w.Code != tc.statusCode {
                t.Errorf("\nExpected:\t%d\nGot:\t\t%d", tc.statusCode, w.Code)
            }
            // Check header
            if ct := w.Header().Get("Content-Type"); ct != "application/json" {
                t.Errorf("\nExpected:\t%s\nGot:\t\t%s", "application/json", ct)
            }
            // Check JSON body
            var resp APIResponse
            err := json.NewDecoder(w.Body).Decode(&resp)
            if err != nil {
                t.Fatalf("Failed to decode response body: %v", err)
            }
            // Check message
            if resp.Message != tc.message {
                t.Errorf("\nExpected:\t%s\nGot:\t\t%s", tc.message, resp.Message)
            }
            // Check error
            if resp.Error != tc.errMsg {
                t.Errorf("\nExpected:\t%s\nGot:\t\t%s", tc.errMsg, resp.Error)
            }
            // Check data
            if !reflect.DeepEqual(resp.Data, tc.data) {
                t.Errorf("\nExpected:\t%v\nGot:\t\t%v", tc.data, resp.Data)
            }
        })
    }

}


