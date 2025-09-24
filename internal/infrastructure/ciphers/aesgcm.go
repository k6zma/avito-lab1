package ciphers

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

type AESGCMCipher struct {
	key  [32]byte
	aead cipher.AEAD
}

func NewAESGCM(key string) (Cipher, error) {
	if len(key) != 32 {
		return nil, ErrInvalidKey
	}

	var k [32]byte

	copy(k[:], key)

	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, fmt.Errorf("failed creating cipher block: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create gcm wrapped cipher block: %w", err)
	}

	return &AESGCMCipher{
		key:  k,
		aead: gcm,
	}, nil
}

func (a *AESGCMCipher) Encrypt(plaintext []byte) ([]byte, error) {
	nonceSize := a.aead.NonceSize()
	nonce := make([]byte, nonceSize)

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed while generating nonce: %w", err)
	}

	ct := a.aead.Seal(nil, nonce, plaintext, nil)
	out := make([]byte, nonceSize+len(ct))

	copy(out, nonce)
	copy(out[nonceSize:], ct)

	return out, nil
}

func (a *AESGCMCipher) Decrypt(data []byte) ([]byte, error) {
	nonceSize := a.aead.NonceSize()
	if len(data) < nonceSize {
		return nil, ErrCorruptedPayload
	}

	nonce := data[:nonceSize]
	ct := data[nonceSize:]

	pt, err := a.aead.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, ErrCorruptedPayload
	}

	return pt, nil
}
