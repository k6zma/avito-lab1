package validators_test

import (
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"

	"github.com/k6zma/avito-lab1/pkg/validators"
)

const (
	capitalizedTestPrefix = "CapitalizedValidator"
)

type capitalizedTestData struct {
	testName   string
	inputValue string
	want       bool
}

type capitalizedValidationStruct struct {
	Value string `validate:"capitalized"`
}

func TestCapitalizedValidator(t *testing.T) {
	tests := []capitalizedTestData{
		{"empty string", "", false},
		{"capitalized", "K6zma", true},
		{"not capitalized", "k6zma", false},
		{"single capitalized letter", "A", true},
		{"single lower letter", "a", false},
		{"multi capitalized string", "Go course", true},
		{"multi lower string", "go course", false},
		{"russian capitalized", "Михаил", true},
		{"russian lower", "михаил", false},
	}

	v := validator.New()
	_ = v.RegisterValidation("capitalized", validators.CapitalizedValidator)

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("[%s]-%s-№%d", capitalizedTestPrefix, tt.testName, i+1),
			func(t *testing.T) {
				s := capitalizedValidationStruct{
					Value: tt.inputValue,
				}

				err := v.Struct(s)

				got := err == nil
				if got != tt.want {
					t.Errorf(
						"capitalized validation failed for input value=%q: got %v, want %v (err: %v)",
						tt.inputValue,
						got,
						tt.want,
						err,
					)
				}
			},
		)
	}
}
