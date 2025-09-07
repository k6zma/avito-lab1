package validators_test

import (
	"fmt"
	"testing"

	"github.com/k6zma/avito-lab1/pkg/validators"
)

const (
	initTestPrefix = "InitValidators"
)

type capitalizedFieldStruct struct {
	ValueCapitalized string `validate:"capitalized"`
}

type initTestCase struct {
	testName   string
	inputValue string
	want       bool
}

func TestInitValidators_Idempotent(t *testing.T) {
	t.Run(fmt.Sprintf("[%s]-idempotentTest-firstCall", initTestPrefix), func(t *testing.T) {
		if err := validators.InitValidators(); err != nil {
			t.Errorf("InitValidators() first call error = %v, want nil", err)
		}
	})

	t.Run(fmt.Sprintf("[%s]-idempotentTest-secondCall", initTestPrefix), func(t *testing.T) {
		if err := validators.InitValidators(); err != nil {
			t.Errorf("InitValidators() second call error = %v, want nil", err)
		}
	})
}

func TestInitValidators_CapitalizedRegistered(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("InitValidators() error = %v, want nil", err)
	}

	tests := []initTestCase{
		{"capitalized", "Test", true},
		{"not capitalized", "test", false},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("[%s]-%s-№%d", initTestPrefix, tt.testName, i+1), func(t *testing.T) {
			s := capitalizedFieldStruct{
				ValueCapitalized: tt.inputValue,
			}

			err := validators.Validate.Struct(s)

			got := err == nil
			if got != tt.want {
				t.Errorf(
					"capitalized validation failed for input value=%q: got %v, want %v (err: %v)",
					tt.inputValue, got, tt.want, err,
				)
			}
		})
	}
}
