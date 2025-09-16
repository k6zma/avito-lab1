package persisters

import "errors"

var ErrMismatchPayloadAndWriteLen = errors.New("mismatch between payload length and write length")
