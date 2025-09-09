package sharedauth

import (
    "testing"
    "errors"
    "encoding/hex"
    "strings"
    "bytes"
)

// python3 -c "import hashlib; print(hashlib.sha256(b'mysalt' + b'mypassword').hexdigest())"
// expected Hash got from python hash
func Test_HashPasswordFn(t *testing.T) {
    salt := "mysalt"
    password := "mypassword"
    expectedHash := "d02878b06efa88579cd84d9e50b211c0a7caa92cf243bad1622c66081f7e2692"
    hashBytes := HashPasswordFn(salt, password)
    hashHex := hex.EncodeToString(hashBytes)
    if expectedHash != hashHex {
        t.Errorf("Hash missmatch:\nExpected:\t%s\nGot:\t\t%s", expectedHash, hashHex)
    }
}


func Test_GenerateRandomBytesFn(t *testing.T) {
    const length = 32

    for i := 0; i < 10; i ++ {
        randomBytes, err := GenerateRandomBytesFn(length)
        if err != nil {
            t.Fatalf("Unedpected error: %v", err)
        }

        if len(randomBytes) != length {
            t.Errorf("Salt wrong length:\nExpected:\t%d\nGot:\t\t%d", length, len(randomBytes))
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
            keyBytes, err := DeriveKeyFn(tc.salt, tc.password, tc.hexLen)
            key := hex.EncodeToString(keyBytes)
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
    cipherBytes, err := EncryptAES(keyBytes, plaintextBytes)
    if err != nil {
        t.Fatalf("Failed to encrypt")
    }

    // Decrypt
    decryptedBytes, err := DecryptAES(keyBytes, cipherBytes)
    if err != nil {
        t.Fatalf("Failed to encrypt")
    }

    if string(decryptedBytes) != "secret" {
        t.Errorf("\nExpected:\t%s\nGot:\t\t%s", "secret", string(decryptedBytes))
    }
}


func Test_WholeProcess(t *testing.T){
    saltConstant := "344feecf40d375380ed5f523b9029647bf7c9f2261e0341a87aa5df6d49c4e31"
    // Generate SYM KEY
    symkeyBytes, err := GenerateRandomBytesFn(32)
    if err != nil {
        t.Fatalf("Failed to generate symetric key")
    }

    // Derive KEY that will be used for encryption
    keyBytes, err := DeriveKeyFn(saltConstant, "strong_password", 32) // hex string
    if err != nil {
        t.Fatalf("Failed to derive key from password")
    }

    // Encrypt
    cipherBytes, err := EncryptAES(keyBytes, symkeyBytes) // double encoding ???
    if err != nil {
        t.Fatalf("Failed to encrypt")
    }

    // Decrypt
    decryptedBytes, err := DecryptAES(keyBytes, cipherBytes)
    if !bytes.Equal(decryptedBytes, symkeyBytes) {
        t.Errorf("Missmatch\nSymKey:\t\t%s\nDecrypted:\t%s",
        string(symkeyBytes), string(decryptedBytes))
    }
}



