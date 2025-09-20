package ciphers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/k6zma/avito-lab1/internal/infrastructure/ciphers"
)

const (
	cipherTestPrefix = "CipherAESGCM"

	testKey = "abcdefghijklmnopqrstuvwxyz123456"
)

type newKeyCase struct {
	name string
	key  string
	ok   bool
}

func TestAESGCM_NewAESGCM_KeyLengthValidation(t *testing.T) {
	tests := []newKeyCase{
		{
			"empty key",
			"",
			false,
		},
		{
			"too short 31 symbols",
			"abcdefghijklmnopqrstuvwxyz12345",
			false,
		},
		{
			"too long 33 symbols",
			"abcdefghijklmnopqrstuvwxyz1234567",
			false,
		},
		{
			"valid (32)",
			testKey,
			true,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("[%s]-new-%s-№%d", cipherTestPrefix, tc.name, i+1), func(t *testing.T) {
			_, err := ciphers.NewAESGCM(tc.key)
			gotOK := err == nil

			if gotOK != tc.ok {
				t.Fatalf("[%s][NewAESGCM] got ok=%v, want ok=%v (err=%v)", cipherTestPrefix, gotOK, tc.ok, err)
			}
		})
	}
}

func TestAESGCM_EncryptDecrypt_RoundTrip(t *testing.T) {
	ctx := context.Background()

	c, err := ciphers.NewAESGCM(testKey)
	if err != nil {
		t.Fatalf("[%s][RoundTrip] failed initing cipher: %v", cipherTestPrefix, err)
	}

	payloads := [][]byte{
		[]byte(""),
		[]byte("hello"),
		[]byte("bye bye"),
		make([]byte, 1024),
	}

	for i, p := range payloads {
		t.Run(fmt.Sprintf("[%s]-roundtrip-case-%d", cipherTestPrefix, i+1), func(t *testing.T) {
			ct, err := c.Encrypt(ctx, p)
			if err != nil {
				t.Fatalf("[%s][Encrypt] unexpected error while encrypting: %v", cipherTestPrefix, err)
			}

			if len(ct) == 0 {
				t.Fatalf("[%s][Encrypt] got empty cipher text", cipherTestPrefix)
			}

			pt, err := c.Decrypt(ctx, ct)
			if err != nil {
				t.Fatalf("[%s][Decrypt] unexpected error while decrypting: %v", cipherTestPrefix, err)
			}

			if string(pt) != string(p) {
				t.Fatalf("[%s][RoundTrip] plaintext mismatch: got=%q want=%q", cipherTestPrefix, pt, p)
			}
		})
	}
}

func TestAESGCM_Decrypt_WithWrongKey_Fails(t *testing.T) {
	ctx := context.Background()

	cipher1, err := ciphers.NewAESGCM(testKey)
	if err != nil {
		t.Fatalf("[%s][WrongKey] failed initing first cipher: %v", cipherTestPrefix, err)
	}

	cipher2, err := ciphers.NewAESGCM("12345678910111213141516171819201")
	if err != nil {
		t.Fatalf("[%s][WrongKey] failed initing second cipher: %v", cipherTestPrefix, err)
	}

	ct, err := cipher1.Encrypt(ctx, []byte("top_secret"))
	if err != nil {
		t.Fatalf("[%s][WrongKey][Encrypt] unexpected error while encrypting: %v", cipherTestPrefix, err)
	}

	if _, err := cipher2.Decrypt(ctx, ct); err == nil {
		t.Fatalf("[%s][WrongKey][Decrypt] expected error while decrypting with wrong key, but got nil", cipherTestPrefix)
	}
}

func TestAESGCM_Decrypt_TamperedCiphertext_Fails(t *testing.T) {
	ctx := context.Background()

	c, err := ciphers.NewAESGCM(testKey)
	if err != nil {
		t.Fatalf("[%s][Tamper] failed initing cipher: %v", cipherTestPrefix, err)
	}

	ct, err := c.Encrypt(ctx, []byte("pau pau pau"))
	if err != nil {
		t.Fatalf("[%s][Tamper][Encrypt] unexpected error while encrypting: %v", cipherTestPrefix, err)
	}

	if len(ct) > 0 {
		ct[0]++
	}

	if _, err := c.Decrypt(ctx, ct); err == nil {
		t.Fatalf("[%s][Tamper][Decrypt] expected error on tampered ciphertext, got nil", cipherTestPrefix)
	}
}
