package validators

import (
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
)

var once sync.Once

var Validate *validator.Validate

func InitValidators() error {
	var initErr error

	once.Do(func() {
		Validate = validator.New()

		initErr = registerStringValidators(Validate)
		if initErr != nil {
			initErr = fmt.Errorf("error while registering string validators: %w", initErr)
		}
	})

	return initErr
}

func registerStringValidators(v *validator.Validate) error {
	err := v.RegisterValidation("capitalized", CapitalizedValidator)
	if err != nil {
		return fmt.Errorf("error while registering capitalized validator: %w", err)
	}

	return nil
}
