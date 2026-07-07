package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

func Encrypt(plaintext, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil { return nil, err }
	gcm, err := cipher.NewGCM(block)
	if err != nil { return nil, err }
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil { return nil, err }
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(ciphertext)))
	base64.StdEncoding.Encode(encoded, ciphertext)
	return encoded, nil
}

func Decrypt(cipherBytes []byte, key string) (string, error) {
	if len(cipherBytes) == 0 { return "", nil }
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(cipherBytes)))
	n, err := base64.StdEncoding.Decode(decoded, cipherBytes)
	if err != nil { return "", err }
	decoded = decoded[:n]
	block, err := aes.NewCipher([]byte(key))
	if err != nil { return "", err }
	gcm, err := cipher.NewGCM(block)
	if err != nil { return "", err }
	nonceSize := gcm.NonceSize()
	if len(decoded) < nonceSize { return "", errors.New("ciphertext too short") }
	nonce, ciphertext := decoded[:nonceSize], decoded[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil { return "", err }
	return string(plaintext), nil
}
