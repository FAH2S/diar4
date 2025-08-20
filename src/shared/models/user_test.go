package sharedmodels
import (
    "testing"
    "strings"
)


// hash and salt for pwd: 'strong_password'
// salt: "344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31"
// hash: "0c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de"


//{{{ validate succ
func Test_UserModel_Validate_Succ(t *testing.T) {
    userData := map[string]string {
        "username":     "valid_test_user1",
        "salt":         "344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
        "hash":         "0c8fd825308df79b313a71b90ee93f7d889207c2277c477b424f83162a5aa4de",
        "encSymkey":    strings.Repeat("0123456789ab", 10),
    }

    // expected outcome
    expected := error(nil)
    // check
    user := User{
        Username:   userData["username"],
        Salt:       userData["salt"],
        Hash:       userData["hash"],
        EncSymkey:  userData["encSymkey"],
    }
    actual := user.Validate()
    if actual != expected {
        t.Errorf("\nExpected: %v\nGot:\n%v", expected, actual)
    }
}
//}}} validate succ


//{{{ username fail
func Test_UserModel_Username_Fail(t *testing.T) {
    tests := []struct {
        name        string
        username    string
        expected    string
    }{
        {
            name:       "UsernameTooShort",
            username:   "ab",
            expected:   "username: length must be between 3 and 30 char long",
        },
        {
            name:       "UsernameTooLong",
            username:   "ab" + strings.Repeat("ab", 30),
            expected:   "username: length must be between 3 and 30 char long",
        },
        {
            name:       "UsernameInvalidChars",
            username:   "invalid_test_user!@#$%^",
            expected:   "username: contains invalid characters",
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            actual := IsValidUsernameFn(tc.username)
            if actual == nil || actual.Error() != tc.expected {
                t.Errorf("\nExpected:\t%q\nGot:\t%q", tc.expected, actual)
            }
        })
    }
}

//}}} username fail


//{{{ hex string fail
func Test_UserModel_HexString_Fail(t *testing.T) {
    tests := []struct {
        name        string
        hexStr      string
        hexStrName  string
        length      int
        expected    string
    }{
        {
            name:       "SaltTooShort",
            hexStr:     "1234567890abcdf",
            hexStrName: "salt",
            length:     64,
            expected:   "salt: length must be exactly 64 char long",
        },{
            name:       "SaltTooLong",
            hexStr:     "" + strings.Repeat("1234567890abcdef", 10),
            hexStrName: "salt",
            length:     64,
            expected:   "salt: length must be exactly 64 char long",
        },{
            name:       "SaltInvalidChars",
            hexStr:     "3T4feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
            hexStrName: "salt",
            length:     64,
            expected:   "salt: contains invalid characters",
        },{
            name:       "HashTooShort",
            hexStr:     "1234567890abcdf",
            hexStrName: "hash",
            length:     64,
            expected:   "hash: length must be exactly 64 char long",
        },{
            name:       "HashTooLong",
            hexStr:     "" + strings.Repeat("1234567890abcdf", 10),
            hexStrName: "hash",
            length:     64,
            expected:   "hash: length must be exactly 64 char long",
        },{
            name:       "hashInvalidChars",
            hexStr:     "3T4feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31",
            hexStrName: "hash",
            length:     64,
            expected:   "hash: contains invalid characters",
        },{
            name:       "encSymkeyTooShort",
            hexStr:     strings.Repeat("0123456789", 5),
            hexStrName: "enc_symkey",
            length:     120,
            expected:   "enc_symkey: length must be exactly 120 char long",
        },{
            name:       "encSymkeyTooLong",
            hexStr:     "" + strings.Repeat("1234567890abcdef", 20),
            hexStrName: "enc_symkey",
            length:     120,
            expected:   "enc_symkey: length must be exactly 120 char long",
        },{
            name:       "hashInvalidChars",
            hexStr:     strings.Repeat("g123456789", 12),
            hexStrName: "enc_symkey",
            length:     120,
            expected:   "enc_symkey: contains invalid characters",
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            actual := IsValidHexStringFn(tc.hexStr, tc.hexStrName, tc.length)
            if actual == nil || actual.Error() != tc.expected {
                t.Errorf("\nExpected:\t%q\nGot:\t%q", tc.expected, actual)
            }
        })
    }
}
//}}} hex string fail


//{{{ Test ValidateUserMap
func Test_ValidateUserMap(t *testing.T) {
    tests := []struct {
        name                string
        input               map[string]interface{}
        expectedErrSubStr   string
    }{
        //field must be string, 
        {
            name:               "FieldNotString",
            input:              map[string]interface{}{
                "salt": 1234,
            },
            expectedErrSubStr:  "field \"salt\" must be string",
        }, {
            // Check if passing err is correct from validate fn
            name:               "FieldNotValid",
            input:              map[string]interface{}{
                "salt": "1234567890abcdf",
            },
            expectedErrSubStr:  "salt: length must be exactly 64 char long",
        },
    }
    // Iterate
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            err := ValidateUserMap(tc.input)
            if err == nil || err.Error() != tc.expectedErrSubStr {
                t.Errorf("\nExpected:\t%q\nGot:\t\t%q", tc.expectedErrSubStr, err)
            }
        })
    }
}
//}}} Test ValidateUserMap

