package persisters

import "errors"

var (
	ErrMismatchPayloadAndWriteLen = errors.New("mismatch between payload length and write length")
	ErrInvalidCipher              = errors.New("invalid cipher provided")
)
