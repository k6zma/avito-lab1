package flags

import (
	"flag"
	"fmt"

	"github.com/k6zma/avito-lab1/pkg/validators"
)

const (
	dataFilePathFlagName         = "data_path"
	dataFilePathFlagDefaultValue = "students_data.json"
	dataFilePathFlagDesc         = "Path to file which will contains data about students"
)

var configPathFlag = flag.String(
	dataFilePathFlagName,
	dataFilePathFlagDefaultValue,
	dataFilePathFlagDesc,
)

type StudyFlags struct {
	ConfigPath string `validate:"required,filepath"`
}

func GetFlags() (*StudyFlags, error) {
	flag.Parse()

	result := &StudyFlags{
		ConfigPath: *configPathFlag,
	}

	if err := validators.Validate.Struct(result); err != nil {
		return nil, fmt.Errorf("error while validating flags in studify app: %w", err)
	}

	return result, nil
}
