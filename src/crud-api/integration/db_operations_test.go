package integration
import (
    "testing"
    "errors"
    "strings"
    "reflect"
    "fmt"
)
import (
    _ "github.com/lib/pq"
    smodels "github.com/FAH2S/diar4/src/shared/models"
    cruduser "github.com/FAH2S/diar4/src/crud-api/user"
)


//{{{ Insert user
func Test_InsertUser(t *testing.T) {
    validSalt :=        "344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31"
    validSaltSymkey :=  "13e4c94d6fd5aa78a1430e1622f9c7fb7469209b325f5de083573d04fceef443"
    validHash :=        "0c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de"
    validEncSymkey :=   "0c8fd08df79b313a71b90ee93f7d889207c2277c477b424f831a5aa4de344feecf40d3753805f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31"
    tests := []struct {
        name                string
        user                smodels.User
        expectedStatusCode  int
        expectedError       error
    }{
        {
            name:               "validInput",
            user:               smodels.User{
                Username:   "test_user_123",
                Salt:       validSalt,
                Hash:       validHash,
                SaltSymkey: validSaltSymkey,
                EncSymkey:  validEncSymkey,
            },
            expectedStatusCode: 201,
            expectedError:      nil,
        }, {
            name:               "userAlreadyExists",
            user:               smodels.User{
                Username:   "test_user_123",
                Salt:       validSalt,
                Hash:       validHash,
                SaltSymkey: validSaltSymkey,
                EncSymkey:  validEncSymkey,
            },
            expectedStatusCode: 409,
            expectedError:      errors.New("user already exists"),
        }, {// Keep in mind this is DB constraint check not .validate (validate is endpoint lvl)
            name:               "unprocessableSalt",
            user:               smodels.User{
                Username:   "test_user_invalid_data",
                Salt:       "invalid_salt",
                Hash:        validHash,
                EncSymkey:  validEncSymkey,
            },
            expectedStatusCode: 422,
            expectedError:      errors.New("violates check constraint \"users_salt_check\""),
        }, {// Keep in mind this is DB constraint check not .validate (validate is endpoint lvl)
            name:               "unprocessableHash",
            user:               smodels.User{
                Username:   "test_user_invalid_data",
                Salt:       validSalt, 
                Hash:       "invalid_hash",
                EncSymkey:  validEncSymkey,
            },
            expectedStatusCode: 422,
            expectedError:      errors.New("violates check constraint \"users_hash_check\""),
        }, {// Keep in mind this is DB constraint check not .validate (validate is endpoint lvl)
            name:               "unprocessableEncSymkey",
            user:               smodels.User{
                Username:   "test_user_invalid_data",
                Salt:       validSalt,
                Hash:       validHash,
                EncSymkey:  "",
            },
            expectedStatusCode: 422,
            expectedError:      errors.New("violates check constraint \"users_enc_symkey_check\""),
        }, {
            name:               "unprocessableSaltEncSymkey",
            user:               smodels.User{
                Username:   "test_user_invalid_data",
                Salt:       validSalt,
                Hash:       validHash,
                SaltSymkey: "invalid_salt",
                EncSymkey:  "",
            },
            expectedStatusCode: 422,
            expectedError:      errors.New("violates check constraint \"users_enc_symkey_check\""),
        }, {
            name:               "unprocessableSaltEncSymkeyNonNullConstraint",
            user:               smodels.User{
                Username:   "test_user_invalid_data",
                Salt:       validSalt,
                Hash:       validHash,
                EncSymkey:  "",
            },
            expectedStatusCode: 422,
            expectedError:      errors.New("violates check constraint \"users_enc_symkey_check\""),
        },
    }

    for _, tc := range tests{
        t.Run(tc.name, func(t *testing.T) {
            actual, err := cruduser.InsertUser(db, tc.user)
            if tc.expectedStatusCode != actual {
                t.Errorf("\nExpected:\t%d\nGot:\t%d", tc.expectedStatusCode, actual)
            }

            if (err == nil) != (tc.expectedError == nil) {
                t.Errorf("\nExpected:\t%v\nGot:\t%v", tc.expectedError, err)
            } else if err != nil &&
                tc.expectedError != nil &&
                !strings.Contains(err.Error(), tc.expectedError.Error()){
                t.Errorf("\nExpected to contain:\t%q\nGot:\t\t\t%q", tc.expectedError, err)
            }
        })
    }
}
//}}} Insert user


//{{{ Select user
func Test_SelectUser(t *testing.T) {
    validSalt :=        "344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31"
    validHash :=        "0c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de"
    validSaltSymkey :=  "13e4c94d6fd5aa78a1430e1622f9c7fb7469209b325f5de083573d04fceef443"
    validEncSymkey :=   "0c8fd08df79b313a71b90ee93f7d889207c2277c477b424f831a5aa4de344feecf40d3753805f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31"
    // Create user that will be fetched (not dependant on previous tests)
    createUser := smodels.User{
        Username:   "test_select_user1",
        Salt:       validSalt,
        Hash:       validHash,
        SaltSymkey: validSaltSymkey,
        EncSymkey:  validEncSymkey,
    }
    _, err := cruduser.InsertUser(db, createUser)
    if err != nil {
        t.Fatalf("Failed to create user that will be read/fetch-ed: %v", err)
    }

    // Test cases and its expected result
    tests := []struct {
        name                string
        username            string
        expectedStatusCode  int
        expectedError       error
        expectedData        *smodels.User
    }{
        {
            name:               "validInput",
            username:           "test_select_user1",
            expectedStatusCode: 200,
            expectedError:      nil,
            expectedData:       &smodels.User{
                Username:   "test_select_user1",
                Salt:       validSalt,
                Hash:       validHash,
                SaltSymkey: validSaltSymkey,
                EncSymkey:  validEncSymkey,
            },
        }, {
            name:               "notFound",
            username:           "not_found",
            expectedStatusCode: 404,
            expectedError:      fmt.Errorf("user not found"),
            expectedData:       nil,
        },

    }
    // Iterate
    for _, tc := range tests{
        t.Run(tc.name, func(t *testing.T) {
            statusCode, user, err := cruduser.SelectUser(db, tc.username)
            if tc.expectedStatusCode != statusCode {
                t.Errorf("\nExpected:\t%d\nGot:\t%d", tc.expectedStatusCode, statusCode)
            }

            if (err == nil) != (tc.expectedError == nil) {
                t.Errorf("\nExpected:\t%v\nGot:\t%v", tc.expectedError, err)
            } else if err != nil &&
                tc.expectedError != nil &&
                !strings.Contains(err.Error(), tc.expectedError.Error()){
                t.Errorf("\nExpected to contain:\t%q\nGot:\t\t\t%q", tc.expectedError, err)
            }
            fmt.Printf("got: %T", user)
            fmt.Printf("want: %T", tc.expectedData)
            if !reflect.DeepEqual(user, tc.expectedData) {
                t.Errorf("\nExpected data not same\nWant:\t%+v\nGot:\t%+v", tc.expectedData, user)
            }
        })
    }
}
//}}} Select user


//{{{ Update user
func Test_UpdateUser(t *testing.T) {
    validSalt :=        "344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31"
    validHash :=        "0c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de"
    validSaltSymkey :=  "13e4c94d6fd5aa78a1430e1622f9c7fb7469209b325f5de083573d04fceef443"
    validEncSymkey :=   "0c8fd08df79b313a71b90ee93f7d889207c2277c477b424f831a5aa4de344feecf40d3753805f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31"
    // Create user that will be updated (not dependant on previous tests)
    createUser := smodels.User{
        Username:   "test_update_user1",
        Salt:       validSalt,
        Hash:       validHash,
        SaltSymkey: validSaltSymkey,
        EncSymkey:  validEncSymkey,
    }
    _, err := cruduser.InsertUser(db, createUser)
    if err != nil {
        t.Fatalf("Failed to create user that will be updated-ed: %v", err)
    }
    tests := []struct {
        name                string
        data                map[string]interface{}
        username            string
        expectedStatusCode  int
        expectedError       error
    }{
        //succ, 422 empty, 404 not found/updated, test that gives non existatnt row
        {
            name:               "Success",
            data:               map[string]interface{}{
                "hash":"1111d825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de",
            },
            username:           "test_update_user1",
            expectedStatusCode: 200,
            expectedError:      nil,
        }, {
            name:               "SuccessSaltSymkey",
            data:               map[string]interface{}{
                "salt_symkey":"1111d825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de",
            },
            username:           "test_update_user1",
            expectedStatusCode: 200,
            expectedError:      nil,
        }, {
            name:               "FailSaltSymkeyConstraint",
            data:               map[string]interface{}{
                "salt_symkey":"111a5aa4de",
            },
            username:           "test_update_user1",
            expectedStatusCode: 422,
            expectedError:      errors.New("violates check constraint"),
        }, {
            name:               "FailNonExistantField400",
            data:               map[string]interface{}{
                "no_hash":"2222d825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de",
            },
            username:           "test_update_user1",
            expectedStatusCode: 400,
            expectedError:      errors.New("unknown column used"),
        }, {
            name:               "FailNoUpdates404",
            data:               map[string]interface{}{
                "hash":"2222d825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de",
            },
            username:           "test_update_user99",
            expectedStatusCode: 404,
            expectedError:      errors.New("no rows were affected"),
        }, {
            name:               "FailDataEmpty422",
            data:               map[string]interface{}{},
            username:           "test_update_user1",
            expectedStatusCode: 422,
            expectedError:      errors.New("no fields to update"),
        },
    }
    //iterate
    for _, tc := range tests{
        t.Run(tc.name, func(t *testing.T) {
            statusCode, err := cruduser.UpdateUser(db, tc.data, tc.username)
            if tc.expectedStatusCode != statusCode {
                t.Errorf("Wrong status code\nExpected:\t%d\nGot:\t\t%d", tc.expectedStatusCode, statusCode)
            }

            if (err == nil) != (tc.expectedError == nil) {
                t.Errorf("Wrong Error\nExpected:\t%v\nGot:\t\t%v", tc.expectedError, err)
            } else if err != nil &&
                tc.expectedError != nil &&
                !strings.Contains(err.Error(), tc.expectedError.Error()){
                t.Errorf("Wrong Error\nExpected:\t%q\nGot:\t\t%q", tc.expectedError, err)
            }
        })
    }
}

//}}} Update user


//{{{ Delete user
func Test_DeleteUser(t *testing.T) {
    validSalt :=        "344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31"
    validHash :=        "0c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de"
    validSaltSymkey :=  "13e4c94d6fd5aa78a1430e1622f9c7fb7469209b325f5de083573d04fceef443"
    validEncSymkey :=   "0c8fd08df79b313a71b90ee93f7d889207c2277c477b424f831a5aa4de344feecf40d3753805f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31"
    // Create user that will be updated (not dependant on previous tests)
    createUser := smodels.User{
        Username:   "test_delete_user1",
        Salt:       validSalt,
        Hash:       validHash,
        SaltSymkey: validSaltSymkey,
        EncSymkey:  validEncSymkey,
    }
    _, err := cruduser.InsertUser(db, createUser)
    if err != nil {
        t.Fatalf("Failed to create user that will be updated-ed: %v", err)
    }
    tests := []struct {
        name                string
        username            string
        expectedStatusCode  int
        expectedError       error
    }{
        {
            name:               "SuccessDeleteUser",
            username:           "test_delete_user1",
            expectedStatusCode: 200,
            expectedError:      nil,
        }, {
            name:               "FailUserNotFound",
            username:           "test_delete_user1",
            expectedStatusCode: 404,
            expectedError:      errors.New("DeleteUser: CheckRowsAffectedFn: no rows were affected"),
        },
    }
    // Iterate
    for _, tc := range tests{
        t.Run(tc.name, func(t *testing.T) {
            statusCode, err := cruduser.DeleteUser(db, tc.username)
            if tc.expectedStatusCode != statusCode {
                t.Errorf("Wrong status code\nExpected:\t%d\nGot:\t\t%d", tc.expectedStatusCode, statusCode)
            }

            if (err == nil) != (tc.expectedError == nil) {
                t.Errorf("Wrong Error\nExpected:\t%v\nGot:\t\t%v", tc.expectedError, err)
            } else if err != nil &&
                tc.expectedError != nil &&
                !strings.Contains(err.Error(), tc.expectedError.Error()){
                t.Errorf("Wrong Error\nExpected:\t%q\nGot:\t\t%q", tc.expectedError, err)
            }
        })
    }
}
//}}} Delete user


