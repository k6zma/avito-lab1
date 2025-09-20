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

	cipherKeyFlagName     = "cipher_key"
	cipherKetDefaultValue = ""
	cipherKeyFlagDesc     = "Key for encryption/decryption of students data using AES-GCM - it's required to be 32 characters long"
)

var configPathFlag = flag.String(
	dataFilePathFlagName,
	dataFilePathFlagDefaultValue,
	dataFilePathFlagDesc,
)

var cipherKeyFlag = flag.String(
	cipherKeyFlagName,
	cipherKetDefaultValue,
	cipherKeyFlagDesc,
)

type StudyFlags struct {
	ConfigPath string `validate:"required,filepath"`
	CipherKey  string `validate:"required,len=32"`
}

func GetFlags() (*StudyFlags, error) {
	flag.Parse()

	result := &StudyFlags{
		ConfigPath: *configPathFlag,
		CipherKey:  *cipherKeyFlag,
	}

	if err := validators.Validate.Struct(result); err != nil {
		return nil, fmt.Errorf("error while validating flags in studify app: %w", err)
	}

	return result, nil
}

func ResetForTests(fs *flag.FlagSet) {
	flag.CommandLine = fs
	
	configPathFlag = flag.String(dataFilePathFlagName, dataFilePathFlagDefaultValue, dataFilePathFlagDesc)
	cipherKeyFlag = flag.String(cipherKeyFlagName, cipherKetDefaultValue, cipherKeyFlagDesc)
}
