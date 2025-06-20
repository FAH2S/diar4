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


func isValidUsernameFn(username string) error {
    // check username, 2 > len > 31, letters + numbers + '_'
    if len(username) < 3 || len(username) > 30 {
        return fmt.Errorf("username: length must be between 3 and 30 char long")
    }
    if !usernameMatch.MatchString(username) {
        return fmt.Errorf("username: contains invalid characters")
    }
    return nil
}


func isValidHexStringFn(hexStr string, hexStrName string, length int) error {
    if len(hexStr) != length {
        return fmt.Errorf("%s: length must be exactly %d char long", hexStrName, length)
    }
    if !hexStrMatch.MatchString(hexStr) {
        return fmt.Errorf("%s: contains invalid characters", hexStrName)
    }
    return nil
}


func (user *User) Validate() error {
    if err := isValidUsernameFn(user.Username); err != nil {
        return err
    }
    if err := isValidHexStringFn(user.Salt, "salt", 64); err != nil {
        return err
    }
    if err := isValidHexStringFn(user.Hash, "hash", 64); err != nil {
        return err
    }
    if err := isValidHexStringFn(user.EncSymkey, "enc_symkey", 120); err != nil {
        return err
    }

    return nil
}
