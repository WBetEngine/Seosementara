package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// EncryptAESGCM encrypts plaintext with a 32-byte key (base64 or raw from env).
func EncryptAESGCM(key []byte, plaintext []byte) (ciphertext, nonce []byte, err error) {
	if len(key) != 32 {
		return nil, nil, errors.New("encryption key must be 32 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}
	nonce = make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}
	return gcm.Seal(nil, nonce, plaintext, nil), nonce, nil
}

func DecryptAESGCM(key []byte, ciphertext, nonce []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("encryption key must be 32 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func KeyFromEnv(b64 string) ([]byte, error) {
	if b64 == "" {
		return nil, errors.New("PIXEL_ENCRYPTION_KEY is required for credential storage")
	}
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	if len(raw) != 32 {
		return nil, errors.New("PIXEL_ENCRYPTION_KEY must decode to 32 bytes")
	}
	return raw, nil
}
