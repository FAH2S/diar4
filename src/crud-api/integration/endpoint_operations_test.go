package integration
import (
    "testing"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "reflect"
    "fmt"
)
import (
    sapi "github.com/FAH2S/diar4/src/shared/api"
    smodels "github.com/FAH2S/diar4/src/shared/models"
    cruduser "github.com/FAH2S/diar4/src/crud-api/user"
    crudmiddleware "github.com/FAH2S/diar4/src/crud-api/middleware"
)

//{{{ DRY
type EndpointTestCase struct {
    Name                string
    Body                string
    ExpectedStatusCode  int
    ExpectedMessage     string
    ExpectedError       string
    ExpectedData        interface{}
}


func assertResponse(
    t *testing.T,
    resp *httptest.ResponseRecorder,
    tc EndpointTestCase,
){
    var bodyResp sapi.APIResponse
    if err := json.Unmarshal(resp.Body.Bytes(), &bodyResp); err != nil {
        t.Fatalf("Failed to parse JSON response as APIResponse: %v", err)
    }
    // status code
    if tc.ExpectedStatusCode != resp.Result().StatusCode {
        t.Errorf("Unexpeted status code:\nGot:\t%d\nWant:\t%d", resp.Result().StatusCode, tc.ExpectedStatusCode)
    }
    // msg
    if !strings.Contains(bodyResp.Message, tc.ExpectedMessage) {
        t.Errorf("Unexpected message:\nGot:\t%s\nWant:\t%s", bodyResp.Message, tc.ExpectedMessage)
    }
    // err
    if !strings.Contains(bodyResp.Error, tc.ExpectedError) {
        t.Errorf("Unexpected error:\nGot:\t%s\nWant:\t%s", bodyResp.Error, tc.ExpectedError)
    }
    // data
    if !reflect.DeepEqual(bodyResp.Data, tc.ExpectedData) {
        t.Errorf("Unexpected data:\nGot:\t%v\nWant:\t%v", bodyResp.Data, tc.ExpectedData)
    }

}
//}}} DRY


//{{{ Middleware
type MiddlewareTestCase struct {
    Name                string
    Method              string
    ContentType         string
    ExpectedStatusCode  int
    ExpectedMessage     string
    ExpectedError       string
    ExpectedData        interface{}
}


// Wrong header
// Wrong method (GET)
func Test_MiddlewareEndpoint_MethodFail(t *testing.T){
    // Mock object, for storing response
    resp := httptest.NewRecorder()
    // Create request
    req := httptest.NewRequest("GET", "/test", nil)
    req.Header.Set("Content-Type", "application/json")

    // Wrap endpoint with middleware
    handler := crudmiddleware.ValidateMethodAndTypeEndpoint(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
    }))
    // Call hander to simulate req and resp
    handler.ServeHTTP(resp, req)
    // Check result
    if resp.Code != 400 {
        t.Errorf("Unexpeted status code:\nGot:\t%d\nWant:\t%d", resp.Code, 400)
    }
}


func Test_MiddlewareEndpoint_HeaderFail(t *testing.T){
    // Mock object, for storing response
    resp := httptest.NewRecorder()
    // Create request
    req := httptest.NewRequest("POST", "/test", nil)
    req.Header.Set("Content-Type", "plain/text")

    // Wrap endpoint with middleware
    handler := crudmiddleware.ValidateMethodAndTypeEndpoint(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
    }))
    // Call hander to simulate req and resp
    handler.ServeHTTP(resp, req)
    // Check result
    if resp.Code != 400 {
        t.Errorf("Unexpeted status code:\nGot:\t%d\nWant:\t%d", resp.Code, 400)
    }
}
//}}} Middleware


//{{{ CreateUserEndpoint
func Test_CreateUserEndpoint(t *testing.T){
    validSalt :=        "344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31"
    validHash :=        "0c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de"
    validEncSymkey :=   "0c8fd08df79b313a71b90ee93f7d889207c2277c477b424f831a5aa4de344feecf40d3753805f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31"
    tests := []EndpointTestCase{
        {
            Name:               "CreateUser",
            Body:               fmt.Sprintf(`{
                "username":"test_user_endpoint1",
                "salt":"%s",
                "hash":"%s",
                "enc_symkey":"%s"
            }`, validSalt, validHash, validEncSymkey),
            ExpectedStatusCode: 201,
            ExpectedMessage:    "Success: create user 'test_user_endpoint1'",
            ExpectedError:      "",
            ExpectedData:       nil,
        }, {
            Name:               "MalformedJSON",
            Body:               fmt.Sprintf(`{
                "username":"test_user_endpoint1",
                "salt":"%s",
                "hash":"0c8f
                `, validSalt),
            ExpectedStatusCode: 400,
            ExpectedMessage:    "Fail: create user ''",
            ExpectedError:      "Invalid JSON",
            ExpectedData:       nil,
        }, {
            Name:               "UserAlreadyExist",
            Body:               fmt.Sprintf(`{
                "username":"test_user_endpoint1",
                "salt":"%s",
                "hash":"%s",
                "enc_symkey":"%s"
            }`, validSalt, validHash, validEncSymkey),
            ExpectedStatusCode: 409,
            ExpectedMessage:    "Fail: create user 'test_user_endpoint1'",
            ExpectedError:      "User already exist",
            ExpectedData:       nil,
        }, {
            Name:               "UnprocessableUsername",
            Body:               fmt.Sprintf(`{
                "username":"fishy user |._.|><|",
                "salt":"%s",
                "hash":"%s",
                "enc_symkey":"%s"
            }`, validSalt, validHash, validEncSymkey),
            ExpectedStatusCode: 422,
            ExpectedMessage:    "Fail: create user 'fishy user |._.|><|'",
            ExpectedError:      "Invalid input format: username: contains invalid characters",
            ExpectedData:       nil,
        }, {
            Name:               "UnprocessableSalt",
            Body:               fmt.Sprintf(`{
                "username":"test_user_endpoint1",
                "salt":"344feecf40d261e0341a87aa5df6d49c4e31",
                "hash":"%s",
                "enc_symkey":"%s"
            }`, validHash, validEncSymkey),
            ExpectedStatusCode: 422,
            ExpectedMessage:    "Fail: create user 'test_user_endpoint1'",
            ExpectedError:      "Invalid input format: salt: length must be exactly 64 char long",
            ExpectedData:       nil,
        }, {
            Name:               "UnprocessableHash",
            Body:               fmt.Sprintf(`{
                "username":"test_user_endpoint1",
                "salt":"%s",
                "hash":"^c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de",
                "enc_symkey":"%s"
            }`, validSalt, validEncSymkey),
            ExpectedStatusCode: 422,
            ExpectedMessage:    "Fail: create user 'test_user_endpoint1'",
            ExpectedError:      "Invalid input format: hash: contains invalid characters",
            ExpectedData:       nil,
        }, {
            Name:               "UnprocessableEncSymkey",
            Body:               fmt.Sprintf(`{
                "username":"test_user_endpoint1",
                "salt":"%s",
                "hash":"%s",
                "enc_symkey":""
            }`, validSalt, validHash),
            ExpectedStatusCode: 422,
            ExpectedMessage:    "Fail: create user 'test_user_endpoint1'",
            ExpectedError:      "Invalid input format: enc_symkey: length must be exactly 120 char long",
            ExpectedData:       nil,
        },
    }

    // Iterate
    for _, tc := range tests {
        t.Run(tc.Name, func(t *testing.T) {
            // Create req, resp
            req := httptest.NewRequest("POST", "/create/user", strings.NewReader(tc.Body))
            resp := httptest.NewRecorder()
            // Call endpoint
            cruduser.CreateUserEndpoint(resp, req, db)
            // Check
            assertResponse(t, resp, tc)
        })
    }
}
//}}} CreateUserEndpoint


//{{{ ReadUserEndpoint
func Test_ReadUserEndpoint(t *testing.T){
    // Create some user that will be fetched
    username := "test_user_read_user1"
    user := smodels.User{
        Username:   username,
        Salt:       "344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
        Hash:       "0c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de",
        EncSymkey:  "0c8fd08df79b313a71b90ee93f7d889207c2277c477b424f831a5aa4de344feecf40d3753805f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
    }
    _, err := cruduser.InsertUser(db, user)
    if err != nil {
        t.Fatalf("Failed to create user that will be read/fetch-ed: %v", err)
    }

    // Define tests and its expected results
    tests := []EndpointTestCase{
        {
            Name:               "ReadUser",
            Body:               fmt.Sprintf(`{"username":"%s"}`, username),
            ExpectedStatusCode: 200,
            ExpectedMessage:    fmt.Sprintf("Success: read user '%s'", username),
            ExpectedError:      "",
            ExpectedData:       map[string]any{
                "username":username,
                "salt":user.Salt,
                "hash":user.Hash,
                "enc_symkey":user.EncSymkey,
            },
        },{
            Name:               "MalformedJSON",
            Body:               fmt.Sprintf(`{"username":"%s`, username),
            ExpectedStatusCode: 400,
            ExpectedMessage:    "Fail: read user ''",
            ExpectedError:      "Invalid JSON",
            ExpectedData:       nil,
        },{
            Name:               "NotFound",
            Body:               `{"username":"not_found"}`,
            ExpectedStatusCode: 404,
            ExpectedMessage:    "Fail: read user 'not_found'",
            ExpectedError:      "User not found, dosen't exist",
            ExpectedData:       nil,
        },{
            Name:               "UnprocessableUsername",
            Body:               `{"username":"fishy user |._.|><|"}`,
            ExpectedStatusCode: 422,
            ExpectedMessage:    "Fail: read user 'fishy user |._.|><|'",
            ExpectedError:      "Invalid input format: username: contains invalid characters",
            ExpectedData:       nil,
        },
    }

    // Iterate
    for _, tc := range tests {
        t.Run(tc.Name, func(t *testing.T) {
            // Create req, resp
            req := httptest.NewRequest("POST", "/read/user", strings.NewReader(tc.Body))
            resp := httptest.NewRecorder()
            // Call endpoint
            cruduser.ReadUserEndpoint(resp, req, db)
            // Check
            assertResponse(t, resp, tc)
        })
    }
}

//}}} ReadUserEndpoint


//{{{ UpdateUserEndpoint
func Test_UpdateUserEndpoint(t *testing.T){
    // Create some user that will be updated 
    username := "test_user_update1"
    user := smodels.User{
        Username:   username,
        Salt:       "344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
        Hash:       "0c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de",
        EncSymkey:  "0c8fd08df79b313a71b90ee93f7d889207c2277c477b424f831a5aa4de344feecf40d3753805f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
    }
    _, err := cruduser.InsertUser(db, user)
    if err != nil {
        t.Fatalf("Failed to create user that will be read/fetch-ed: %v", err)
    }
    tests := []EndpointTestCase{
        {
            Name:               "UpdateUser",
            Body:               `{
                "username":"test_user_update1",
                "salt":"111feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31"
            }`,
            ExpectedStatusCode: 200,
            ExpectedMessage:    "Success: update user 'test_user_update1'",
            ExpectedError:      "",
            ExpectedData:       map[string]any{
                "username":"test_user_update1",
            },
        }, {
            Name:               "MalformedJson",
            Body:               fmt.Sprintf(`{
                "username":"test_user_update1",
                "salt":"111feecf4
            `),
            ExpectedStatusCode: 400,
            ExpectedMessage:    "Fail: update user ''",
            ExpectedError:      "Invalid JSON",
            ExpectedData:       nil,
        }, {
            Name:               "UpdateUserNotFound",
            Body:               fmt.Sprintf(`{
                "username":"test_user_update",
                "salt":"111feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31"
            }`),
            ExpectedStatusCode: 404,
            ExpectedMessage:    "Fail: update user 'test_user_update'",
            ExpectedError:      "User not found, dosen't exist",
            ExpectedData:       nil,
        }, {
            Name:               "InvalidInputFormatSalt",
            Body:               fmt.Sprintf(`{
                "username":"test_user_update1",
                "salt":"111feecf"
            }`),
            ExpectedStatusCode: 422,
            ExpectedMessage:    "Fail: update user ''",
            ExpectedError:      "Invalid input format: salt",
            ExpectedData:       nil,
        }, {
            Name:               "InvalidInputFormatHash",
            Body:               fmt.Sprintf(`{
                "username":"test_user_update1",
                "hash":"111feecf"
            }`),
            ExpectedStatusCode: 422,
            ExpectedMessage:    "Fail: update user ''",
            ExpectedError:      "Invalid input format: hash",
            ExpectedData:       nil,
        }, {
            Name:               "InvalidInputFormatEncSymkey",
            Body:               fmt.Sprintf(`{
                "username":"test_user_update1",
                "enc_symkey":"111feecf"
            }`),
            ExpectedStatusCode: 422,
            ExpectedMessage:    "Fail: update user ''",
            ExpectedError:      "Invalid input format: enc_symkey",
            ExpectedData:       nil,
        }, {
            Name:               "MissingUsername",
            Body:               fmt.Sprintf(`{
                "salt":"111feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31"
            }`),
            ExpectedStatusCode: 422,
            ExpectedMessage:    "Fail: update user ''",
            ExpectedError:      "Missing required field: 'username'",
            ExpectedData:       nil,
        }, {
            Name:               "InvalidInputNotEnoughFileds",
            Body:               fmt.Sprintf(`{
                "username":"test_user_update1"
            }`),
            ExpectedStatusCode: 422,
            ExpectedMessage:    "Fail: update user ''",
            ExpectedError:      "Invalid input: must contain at least 2 fields total",
            ExpectedData:       nil,
        },

    }
    // Iterate
    for _, tc := range tests {
        t.Run(tc.Name, func(t *testing.T) {
            // Create req, resp
            req := httptest.NewRequest("POST", "/update/user", strings.NewReader(tc.Body))
            resp := httptest.NewRecorder()
            // Call endpoint
            cruduser.UpdateUserEndpoint(resp, req, db)
            // Check
            assertResponse(t, resp, tc)
        })
    }
}

//}}} UpdateUserEndpoint




