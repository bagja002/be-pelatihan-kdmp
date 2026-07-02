// Package crypto provides AES-256-GCM authenticated encryption used for
// field-level encryption of sensitive data at rest.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"sync"
)

var (
	key   []byte
	keyMu sync.RWMutex
)

// ErrKeyNotSet is returned when Encrypt/Decrypt are called before SetKey.
var ErrKeyNotSet = errors.New("crypto: encryption key not set")

// SetKey installs the 32-byte AES-256 key. Call once at startup.
func SetKey(k []byte) error {
	if len(k) != 32 {
		return errors.New("crypto: key must be exactly 32 bytes for AES-256")
	}
	keyMu.Lock()
	defer keyMu.Unlock()
	key = append([]byte(nil), k...)
	return nil
}

func currentKey() ([]byte, error) {
	keyMu.RLock()
	defer keyMu.RUnlock()
	if len(key) != 32 {
		return nil, ErrKeyNotSet
	}
	return key, nil
}

// Encrypt returns base64(nonce||ciphertext) of plaintext using AES-256-GCM.
// GCM provides confidentiality and integrity (tamper detection).
func Encrypt(plaintext string) (string, error) {
	k, err := currentKey()
	if err != nil {
		return "", err
	}
	gcm, err := newGCM(k)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

// Decrypt reverses Encrypt. It fails if the data was tampered with.
func Decrypt(encoded string) (string, error) {
	k, err := currentKey()
	if err != nil {
		return "", err
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	gcm, err := newGCM(k)
	if err != nil {
		return "", err
	}
	ns := gcm.NonceSize()
	if len(data) < ns {
		return "", errors.New("crypto: ciphertext too short")
	}
	nonce, ciphertext := data[:ns], data[ns:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func newGCM(k []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}
