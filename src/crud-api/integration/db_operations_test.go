package integration
import (
    "testing"
    "errors"
    "strings"
)
import (
    _ "github.com/lib/pq"
    smodels "github.com/FAH2S/diar4/src/shared/models"
    cruduser "github.com/FAH2S/diar4/src/crud-api/user"
)


//{{{ Insert user
func TestInsertUser(t *testing.T) {
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
                Salt:       "344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
                Hash:       "0c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de",
                EncSymkey:  "0c8fd08df79b313a71b90ee93f7d889207c2277c477b424f831a5aa4de344feecf40d3753805f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
            },
            expectedStatusCode: 201,
            expectedError:      nil,
        }, {
            name:               "userAlreadyExists",
            user:               smodels.User{
                Username:   "test_user_123",
                Salt:       "344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
                Hash:       "0c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de",
                EncSymkey:  "0c8fd08df79b313a71b90ee93f7d889207c2277c477b424f831a5aa4de344feecf40d3753805f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
            },
            expectedStatusCode: 409,
            expectedError:      errors.New("user already exists"),
        }, {// Keep in mind this is DB constraint check not .validate (validate is endpoint lvl)
            name:               "unprocessableSalt",
            user:               smodels.User{
                Username:   "test_user_invalid_data",
                Salt:       "344feecf40d375380ed5fe0341a87aa5df6d49c4e31",
                Hash:       "0c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de",
                EncSymkey:  "0c8fd08df79b313a71b90ee93f7d889207c2277c477b424f831a5aa4de344feecf40d3753805f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
            },
            expectedStatusCode: 422,
            expectedError:      errors.New("violates check constraint \"users_salt_check\""),
        }, {// Keep in mind this is DB constraint check not .validate (validate is endpoint lvl)
            name:               "unprocessableHash",
            user:               smodels.User{
                Username:   "test_user_invalid_data",
                Salt:       "344feecf40d375380ed5fe0341a87aa5df6d49c4e31",
                Hash:       "0c8fd825308df79b313a7$b90ee93f7d889207c2277c477b424f83162a5aa4de",
                EncSymkey:  "0c8fd08df79b313a71b90ee93f7d889207c2277c477b424f831a5aa4de344feecf40d3753805f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
            },
            expectedStatusCode: 422,
            expectedError:      errors.New("violates check constraint \"users_hash_check\""),
        }, {// Keep in mind this is DB constraint check not .validate (validate is endpoint lvl)
            name:               "unprocessableEncSymkey",
            user:               smodels.User{
                Username:   "test_user_invalid_data",
                Salt:       "344feecf40d375380ed5fe0341a87aa5df6d49c4e31",
                Hash:       "0c8fd825308df79b313a7$b90ee93f7d889207c2277c477b424f83162a5aa4de",
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
func TestSelectUser(t *testing.T) {
    tests := []struct {
        name                string
        username            string
        expectedStatusCode  int
        expectedError       error
        expectedData        interface{}
    }{
        {
            name:               "validInput",
            username:           "test_user_123",
            expectedStatusCode: 200,
            expectedError:      nil,
            expectedData:       smodels.User{
                Username:   "test_user_123",
                Salt:       "344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
                Hash:       "0c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de",
                EncSymkey:  "0c8fd08df79b313a71b90ee93f7d889207c2277c477b424f831a5aa4de344feecf40d3753805f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
            },
        }, {
            name:               "notFound",
            username:           "not_found",
            expectedStatusCode: 404,
            expectedError:      nil,
            expectedData:       smodels.User{},
        },

    }
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
            if tc.expectedData != user {
                t.Errorf("\nExpected data not same\nWant:\t%+v\nGot:\t%+v", tc.expectedData, user)
            }
        })
    }
}
//}}} Select user



