package encrypt

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "io"
)

// AESEncryptGCM 使用 AES-256-GCM 加密并返回 base64 文本
func AESEncryptGCM(plaintext, key string) (string, error) {
    if len(key) != 32 { return "", errors.New("AES key must be 32 bytes") }
    block, err := aes.NewCipher([]byte(key))
    if err != nil { return "", err }
    gcm, err := cipher.NewGCM(block)
    if err != nil { return "", err }
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil { return "", err }
    sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    return base64.StdEncoding.EncodeToString(sealed), nil
}

// AESDecryptGCM 解密 base64 文本
func AESDecryptGCM(ciphertextB64, key string) (string, error) {
    if len(key) != 32 { return "", errors.New("AES key must be 32 bytes") }
    data, err := base64.StdEncoding.DecodeString(ciphertextB64)
    if err != nil { return "", err }
    block, err := aes.NewCipher([]byte(key))
    if err != nil { return "", err }
    gcm, err := cipher.NewGCM(block)
    if err != nil { return "", err }
    if len(data) < gcm.NonceSize() { return "", errors.New("ciphertext too short") }
    nonce := data[:gcm.NonceSize()]
    ct := data[gcm.NonceSize():]
    pt, err := gcm.Open(nil, nonce, ct, nil)
    if err != nil { return "", err }
    return string(pt), nil
}

