package sharedauth

import (
    "testing"
    "regexp"
    "errors"
    "encoding/hex"
    "strings"
)


var (
    hexRegex = regexp.MustCompile(`^[0-9a-fA-F]+$`)
)

func Test_GenerateHexStrFn(t *testing.T) {
    const expectedLen = 64

    for i := 0; i < 10; i ++ {
        hexStr, err := GenerateHexStrFn(expectedLen)
        if err != nil {
            t.Fatalf("Unedpected error: %v", err)
        }

        if len(hexStr) != expectedLen {
            t.Errorf("Salt wrong length:\nExpected:\t%d\nGot:\t\t%d", expectedLen, len(hexStr))
        }

        if !hexRegex.MatchString(hexStr) {
            t.Errorf("Salt contains non-hex characters: %s", hexStr)
        }
    }
}


/* Python code used to derive key
python3 -c "import hashlib, binascii; print(
    binascii.hexlify(
        hashlib.pbkdf2_hmac(
            'sha256',
            b'strong_password123',
            binascii.unhexlify('344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31'),
            100000,
            32
        )
    ).decode()
)"
'c59732f89cb8f1daf594d946089d553b2b64c6170e5d69e7cc1e0815dc7c94b2'
*/
func Test_DeriveKeyFn(t *testing.T) {
    saltConstant := "344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31"
    tests := []struct {
        name        string
        salt        string
        password    string
        hexLen      int
        expectedStr string
        expectedErr error
    }{
        {
            name:           "FailHexLengthOdd",
            salt:           "salt_string",
            password:       "password_string",
            hexLen:         33,
            expectedStr:    "",
            expectedErr:    errors.New("hex length must be even,"),
        }, {
            name:           "SuccessVsGODerivedKey",
            salt:           saltConstant,
            password:       "strong_password",
            hexLen:         64,
            expectedStr:    "725cee33fe79c8b6cda5065c8b180f947d7c34d4ee846645ec233ece6c49b019",
            expectedErr:    nil,
        }, {
            name:           "SuccessVsPythonDerivedKey",
            salt:           saltConstant,
            password:       "strong_password123",
            hexLen:         64,
            expectedStr:    "c59732f89cb8f1daf594d946089d553b2b64c6170e5d69e7cc1e0815dc7c94b2",
            expectedErr:    nil,
        },
    }
    // Iterate
    for _, tc := range tests{
        t.Run(tc.name, func(t *testing.T) {
            key, err := DeriveKeyFn(tc.salt, tc.password, tc.hexLen)
            if tc.expectedStr != key {
                t.Errorf("\nExpected:\t%s\nGot:\t\t%s", tc.expectedStr, key)
            }
            if (err == nil) != (tc.expectedErr == nil) {
                t.Errorf("\nExpected:\t%v\nGot:\t\t%v", tc.expectedErr, err)
            } else if err != nil &&
                tc.expectedErr != nil &&
                !strings.Contains(err.Error(), tc.expectedErr.Error()){
                t.Errorf("\nExpected:\t%q\nGot:\t\t%q", tc.expectedErr, err)
            }
        })
    }
}



func Test_EncryptAES(t *testing.T){
    keyBytes, _ := hex.DecodeString("725cee33fe79c8b6cda5065c8b180f947d7c34d4ee846645ec233ece6c49b019")
    plaintextBytes := []byte("secret")

    // Encrypt
    cipher, err := EncryptAES(keyBytes, plaintextBytes)
    if err != nil {
        t.Fatalf("Failed to encrypt")
    }

    // Decrypt
    cipherBytes, _ := hex.DecodeString(cipher)
    plaintext, err := DecryptAES(keyBytes, cipherBytes)
    if err != nil {
        t.Fatalf("Failed to encrypt")
    }

    if plaintext != "secret" {
        t.Errorf("\nExpected:\t%s\nGot:\t\t%s", "secret", plaintext)
    }
}






