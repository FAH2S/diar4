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
func HashPasswordFn(salt, password string) []byte {
    hash := sha256.New()
    hash.Write([]byte(salt))
    hash.Write([]byte(password))
    return hash.Sum(nil)
}


func GenerateRandomBytesFn(length int) ([]byte, error) {
    // 32 bytes = 256 bits = 64 hex char
    fn := "GenerateRandomBytesFn"
    bytes := make([]byte, length)
    _, err := rand.Read(bytes)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", fn, err)
    }

    return bytes, nil
}


func DeriveKeyFn(salt, password string, keyLen int) ([]byte, error) {
    fn := "DeriveKeyFn"
    // PBKDF2 conf
    iterations := 100_000

    saltBytes, err := hex.DecodeString(salt)
    if err != nil {
        return nil, fmt.Errorf("%s: %v", fn, err)
    }

    key := pbkdf2.Key([]byte(password), saltBytes, iterations, keyLen, sha256.New)
    return key, nil
}


func EncryptAES(keyBytes, plaintextBytes []byte) ([]byte, error) {
    wrap := "EncryptAES"
    // AES cipher block
    block, err := aes.NewCipher(keyBytes)
    if err != nil {
        return nil, fmt.Errorf("%s: %v", wrap, err)
    }

    // GCM mode
    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("%s: %v", wrap, err)
    }

    // Nonce
    nonce := make([]byte, aesGCM.NonceSize())
    rand.Read(nonce)

    // Encrypt
    ciphertext := aesGCM.Seal(nil, nonce, plaintextBytes, nil)
    final := append(nonce, ciphertext...)
    return final, nil
}


func DecryptAES(keyBytes, cipherBytes []byte) ([]byte, error) {
    wrap := "DecryptAES"
    // AES cipher block
    block, err := aes.NewCipher(keyBytes)
    if err != nil {
        return nil, fmt.Errorf("%s: %v", wrap, err)
    }

    // GCM mode
    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("%s: %v", wrap, err)
    }

    nonceSize := aesGCM.NonceSize()
    if len(cipherBytes) < nonceSize {
        return nil, fmt.Errorf("%s: ciphertext too short", wrap)
    }
    // cipher = nonce + actual ciphertext
    nonce := cipherBytes[:nonceSize]
    ciphertextBytes := cipherBytes[nonceSize:]

    decryptedBytes, err := aesGCM.Open(nil, nonce, ciphertextBytes, nil)
    if err != nil {
        return nil, fmt.Errorf("%s: %v", wrap, err)
    }
    return decryptedBytes, nil
}




