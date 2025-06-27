package integration
import (
    "testing"
    "database/sql"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "reflect"
)
import (
    sapi "github.com/FAH2S/diar4/src/shared/api"
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
func TestMiddlewareEndpointMethodFail(t *testing.T){
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


func TestMiddlewareEndpointHeaderFail(t *testing.T){
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
func TestCreateUserEndpoint(t *testing.T){
    tests := []EndpointTestCase{
        {
            Name:               "CreateUser",
            Body:               `{
                "username":"test_user_endpoint1",
                "salt":"344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
                "hash":"0c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de",
                "enc_symkey":"0c8fd08df79b313a71b90ee93f7d889207c2277c477b424f831a5aa4de344feecf40d3753805f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31" }`,
            ExpectedStatusCode: 201,
            ExpectedMessage:    "Success: create user 'test_user_endpoint1'",
            ExpectedError:      "",
            ExpectedData:       nil,
        }, {
            Name:               "MalformedJSON",
            Body:               `{
                "username":"test_user_endpoint1",
                "salt":"344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
                "hash":"0c8fd825308df79b313a71b90ee93f7d889207c2277c47`,
            ExpectedStatusCode: 400,
            ExpectedMessage:    "Fail: create user ''",
            ExpectedError:      "Invalid JSON",
            ExpectedData:       nil,
        }, {
            Name:               "UserAlreadyExist",
            Body:               `{
                "username":"test_user_endpoint1",
                "salt":"344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
                "hash":"0c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de",
                "enc_symkey":"0c8fd08df79b313a71b90ee93f7d889207c2277c477b424f831a5aa4de344feecf40d3753805f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31"
            }`,
            ExpectedStatusCode: 409,
            ExpectedMessage:    "Fail: create user 'test_user_endpoint1'",
            ExpectedError:      "User already exist",
            ExpectedData:       nil,
        }, {
            Name:               "UnprocessableUsername",
            Body:               `{
                "username":"fishy user |._.|><|",
                "salt":"344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
                "hash":"0c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de",
                "enc_symkey":"0c8fd08df79b313a71b90ee93f7d889207c2277c477b424f831a5aa4de344feecf40d3753805f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31"
            }`,
            ExpectedStatusCode: 422,
            ExpectedMessage:    "Fail: create user 'fishy user |._.|><|'",
            ExpectedError:      "Invalid input format: username: contains invalid characters",
            ExpectedData:       nil,
        }, {
            Name:               "UnprocessableSalt",
            Body:               `{
                "username":"test_user_endpoint1",
                "salt":"344feecf40d261e0341a87aa5df6d49c4e31",
                "hash":"0c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de",
                "enc_symkey":"0c8fd08df79b313a71b90ee93f7d889207c2277c477b424f831a5aa4de344feecf40d3753805f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31"
            }`,
            ExpectedStatusCode: 422,
            ExpectedMessage:    "Fail: create user 'test_user_endpoint1'",
            ExpectedError:      "Invalid input format: salt: length must be exactly 64 char long",
            ExpectedData:       nil,
        }, {
            Name:               "UnprocessableHash",
            Body:               `{
                "username":"test_user_endpoint1",
                "salt":"344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
                "hash":"^c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de",
                "enc_symkey":"0c8fd08df79b313a71b90ee93f7d889207c2277c477b424f831a5aa4de344feecf40d3753805f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31"
            }`,
            ExpectedStatusCode: 422,
            ExpectedMessage:    "Fail: create user 'test_user_endpoint1'",
            ExpectedError:      "Invalid input format: hash: contains invalid characters",
            ExpectedData:       nil,
        }, {
            Name:               "UnprocessableEncSymkey",
            Body:               `{
                "username":"test_user_endpoint1",
                "salt":"344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
                "hash":"0c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de",
                "enc_symkey":""
            }`,
            ExpectedStatusCode: 422,
            ExpectedMessage:    "Fail: create user 'test_user_endpoint1'",
            ExpectedError:      "Invalid input format: enc_symkey: length must be exactly 120 char long",
            ExpectedData:       nil,
        },
    }

    connStr, err := getPostgresConnStr()
    if err != nil {
        t.Fatalf("Failed to get conn string: %v", err)
    }
    db, err := sql.Open("postgres", connStr)//db
    if err != nil {
        t.Fatalf("Failed to open DB: %v", err)
    }
    defer db.Close()
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


