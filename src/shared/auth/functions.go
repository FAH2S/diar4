package sharedauth

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
)

import (
    "golang.org/x/crypto/pbkdf2"
)

// generate salt, generate symkey, generate/derive key form pwd,
// hash(salt + pwd), encrypt (symkey),

// TODO: doc it
func GenerateHexStrFn(hexLen int) (string, error) {
    // 32 bytes = 256 bits = 64 hex char
    fn := "GenerateHexStrFn"
    if hexLen%2 != 0 {
        return "", fmt.Errorf("%s: hex length must be even, got %d", fn, hexLen)
    }

    bytes := make([]byte, hexLen/2)
    _, err := rand.Read(bytes)
    if err != nil {
        return "", fmt.Errorf("%s: %w", fn, err)
    }
    // Convert to hex string
    return hex.EncodeToString(bytes), nil
}


// TODO: doc it
func DeriveKeyFn(salt, password string, hexLen int) (string, error) {
    fn := "DeriveKeyFn"
    if hexLen%2 != 0 {
        return "", fmt.Errorf("%s: hex length must be even, got %d", fn, hexLen)
    }
    // PBKDF2 conf
    iterations := 100_000
    keyLen := hexLen/2

    saltBytes, err := hex.DecodeString(salt)
    if err != nil {
        return "", fmt.Errorf("%s: %v", fn, err)
    }

    key := pbkdf2.Key([]byte(password), saltBytes, iterations, keyLen, sha256.New)
    return hex.EncodeToString(key), nil
}


// TODO: doc it, maybe FN and not wrapper
func EncryptAES(keyBytes, plaintextBytes []byte) (string, error) {
    wrap := "EncryptAES"
    // AES cipher block
    block, err := aes.NewCipher(keyBytes)
    if err != nil {
        return "", fmt.Errorf("%s: %v", wrap, err)
    }

    // GCM mode
    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return "", fmt.Errorf("%s: %v", wrap, err)
    }

    // Nonce
    nonce := make([]byte, aesGCM.NonceSize())
    rand.Read(nonce)

    // Encrypt
    ciphertext := aesGCM.Seal(nil, nonce, plaintextBytes, nil)
    final := append(nonce, ciphertext...)
    return hex.EncodeToString(final), nil
}
func DecryptAES(keyBytes, cipherBytes []byte) (string, error) {
    wrap := "DecryptAES"
    // AES cipher block
    block, err := aes.NewCipher(keyBytes)
    if err != nil {
        return "", fmt.Errorf("%s: %v", wrap, err)
    }

    // GCM mode
    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return "", fmt.Errorf("%s: %v", wrap, err)
    }

    nonceSize := aesGCM.NonceSize()
    if len(cipherBytes) < nonceSize {
        return "", fmt.Errorf("%s: ciphertext too short", wrap)
    }
    // cipher = nonce + actual ciphertext
    nonce := cipherBytes[:nonceSize]
    ciphertextBytes := cipherBytes[nonceSize:]

    plaintext, err := aesGCM.Open(nil, nonce, ciphertextBytes, nil)
    if err != nil {
        return "", fmt.Errorf("%s: %v", wrap, err)
    }
    return string(plaintext), nil
}



// TODO: test if enc dec works


