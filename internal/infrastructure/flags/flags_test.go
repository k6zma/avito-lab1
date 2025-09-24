package flags_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/k6zma/avito-lab1/internal/infrastructure/flags"
	"github.com/k6zma/avito-lab1/pkg/validators"
)

const (
	flagsTestPrefix  = "StudifyFlags"
	studentsDataPath = "data/students.json"
	cipherKey        = "abcdefghijklmnopqrstuvwxyz123456"
)

type getFlagsCase struct {
	name     string
	dataPath string
	key      string
	wantErr  bool
}

func TestGetFlags_TableDriven(t *testing.T) {
	if validators.Validate == nil {
		if err := validators.InitValidators(); err != nil {
			t.Fatalf("[%s][InitValidators] failed to init validators: %v", flagsTestPrefix, err)
		}
	}

	tests := []getFlagsCase{
		{
			name:     "ok with defaults (only key provided)",
			dataPath: studentsDataPath,
			key:      cipherKey,
			wantErr:  false,
		},
		{
			name:     "ok with custom config path",
			dataPath: "custom/students.json",
			key:      cipherKey,
			wantErr:  false,
		},
		{
			name:     "invalid because empty key)",
			dataPath: studentsDataPath,
			key:      "",
			wantErr:  true,
		},
		{
			name:     "invalid because key too short",
			dataPath: studentsDataPath,
			key:      "short_key",
			wantErr:  true,
		},
		{
			name:     "invalid because key too long",
			dataPath: studentsDataPath,
			key:      cipherKey + "pupu",
			wantErr:  true,
		},
		{
			name:     "invalid because empty path",
			dataPath: "",
			key:      cipherKey,
			wantErr:  true,
		},
	}

	origArgs := os.Args

	defer func() {
		os.Args = origArgs
	}()

	for i, tc := range tests {
		t.Run(
			fmt.Sprintf("[%s]-GetFlags-%s-№%d", flagsTestPrefix, tc.name, i+1),
			func(t *testing.T) {
				flags.ResetForTests(flag.NewFlagSet("studify", flag.ContinueOnError))

				args := []string{"studify"}

				if tc.dataPath != "" {
					args = append(args, fmt.Sprintf("-%s=%s", "data_path", tc.dataPath))
				} else {
					args = append(args, fmt.Sprintf("-%s=%s", "data_path", ""))
				}

				if tc.key != "" {
					args = append(args, fmt.Sprintf("-%s=%s", "cipher_key", tc.key))
				} else {
					args = append(args, fmt.Sprintf("-%s=%s", "cipher_key", ""))
				}

				os.Args = args

				got, err := flags.GetFlags()
				gotErr := err != nil

				if gotErr != tc.wantErr {
					t.Fatalf(
						"[%s][GetFlags] got error=%v, want error=%v (err=%v)",
						flagsTestPrefix, gotErr, tc.wantErr, err,
					)
				}

				if !tc.wantErr {
					if got.ConfigPath != tc.dataPath {
						t.Fatalf(
							"[%s][GetFlags] Data Path mismatch: got=%q want=%q",
							flagsTestPrefix, got.ConfigPath, tc.dataPath,
						)
					}

					if got.CipherKey != tc.key {
						t.Fatalf(
							"[%s][GetFlags] Cipher Key mismatch: got=%q want=%q",
							flagsTestPrefix, got.CipherKey, tc.key,
						)
					}
				}
			},
		)
	}
}

func TestGetFlags_IdempotentValues(t *testing.T) {
	if validators.Validate == nil {
		if err := validators.InitValidators(); err != nil {
			t.Fatalf("[%s][InitValidators] failed to init validators: %v", flagsTestPrefix, err)
		}
	}

	origArgs := os.Args

	defer func() {
		os.Args = origArgs
	}()

	argv := []string{
		"studify",
		fmt.Sprintf("-%s=%s", "data_path", studentsDataPath),
		fmt.Sprintf("-%s=%s", "cipher_key", cipherKey),
	}

	flags.ResetForTests(flag.NewFlagSet("studify", flag.ContinueOnError))

	os.Args = append([]string(nil), argv...)

	first, err := flags.GetFlags()
	if err != nil {
		t.Fatalf(
			"[%s][Idempotent-first] unexpected error while getting flags: %v",
			flagsTestPrefix,
			err,
		)
	}

	flags.ResetForTests(flag.NewFlagSet("studify", flag.ContinueOnError))

	os.Args = append([]string(nil), argv...)

	second, err := flags.GetFlags()
	if err != nil {
		t.Fatalf(
			"[%s][Idempotent-second] unexpected error while getting flags: %v",
			flagsTestPrefix,
			err,
		)
	}

	if first.ConfigPath != second.ConfigPath || first.CipherKey != second.CipherKey {
		t.Fatalf(
			"[%s][Idempotent] mismatch between GetFlags calls: first={%q,%q} second={%q,%q}",
			flagsTestPrefix,
			first.ConfigPath, first.CipherKey,
			second.ConfigPath, second.CipherKey,
		)
	}
}
