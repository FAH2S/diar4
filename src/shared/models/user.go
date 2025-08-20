package sharedmodels
import (
    "fmt"
    "regexp"
)


// Regex init
var (
    usernameMatch = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
    hexStrMatch = regexp.MustCompile(`^[0-9a-fA-F]+$`)
)


// Create user struct
type User struct {
    Username    string `json:"username"`
    Salt        string `json:"salt"`
    Hash        string `json:"hash"`
    EncSymkey   string `json:"enc_symkey"`
}


func IsValidUsernameFn(username string) error {
    // check username, 2 > len > 31, letters + numbers + '_'
    if len(username) < 3 || len(username) > 30 {
        return fmt.Errorf("username: length must be between 3 and 30 char long")
    }
    if !usernameMatch.MatchString(username) {
        return fmt.Errorf("username: contains invalid characters")
    }
    return nil
}


func IsValidHexStringFn(hexStr string, hexStrName string, length int) error {
    if len(hexStr) != length {
        return fmt.Errorf("%s: length must be exactly %d char long", hexStrName, length)
    }
    if !hexStrMatch.MatchString(hexStr) {
        return fmt.Errorf("%s: contains invalid characters", hexStrName)
    }
    return nil
}


func (user *User) Validate() error {
    if err := IsValidUsernameFn(user.Username); err != nil {
        return err
    }
    if err := IsValidHexStringFn(user.Salt, "salt", 64); err != nil {
        return err
    }
    if err := IsValidHexStringFn(user.Hash, "hash", 64); err != nil {
        return err
    }
    if err := IsValidHexStringFn(user.EncSymkey, "enc_symkey", 120); err != nil {
        return err
    }

    return nil
}


func ValidateUserMap(input map[string]interface{}) error {
    // Define validators
    validators := map[string]func(string) error {
        "username": func(val string) error {
            return IsValidUsernameFn(val)
        },
        "salt": func(val string) error {
            return IsValidHexStringFn(val, "salt", 64)
        },
        "hash": func(val string) error {
            return IsValidHexStringFn(val, "hash", 64)
        },
        "enc_symkey": func(val string) error {
            return IsValidHexStringFn(val, "enc_symkey", 120)
        },
    }

    // Iterate
    for field, validateFn := range validators {
        rawInputField, ok := input[field]
        if !ok {
            continue // Skip
        }
        // String check
        strVal, ok := rawInputField.(string)
        if !ok {
            return fmt.Errorf("field %q must be string", field)
        }
        // Validate
        if err := validateFn(strVal); err != nil {
            return err
        }
    }
    return nil
}



