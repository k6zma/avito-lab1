package validators

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

func CapitalizedValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	runes := []rune(value)
	if len(runes) == 0 {
		return false
	}

	return unicode.IsUpper(runes[0])
}
