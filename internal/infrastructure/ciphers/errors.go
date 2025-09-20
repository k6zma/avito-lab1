package ciphers

import "errors"

var (
	ErrInvalidKey       = errors.New("invalid encryption key")
	ErrEncryptFailed    = errors.New("encrypt failed")
	ErrCorruptedPayload = errors.New("ciphertext/auth tag is invalid or payload is corrupted")
)
