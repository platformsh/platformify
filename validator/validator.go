package validator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v3"
)

// ValidateFile checks the file exists and is valid yaml, then returns the unmarshalled data.
func ValidateFile(path string) (interface{}, error) {
	// Does the file exist?
	yamlData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read .platform.app.yaml file: %v", err)
	}

	// Does the yaml unmarshal properly?
	var data interface{}
	err = yaml.Unmarshal(yamlData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML data: %v", err)
	}

	// Return unmarshalled data.
	return data, nil
}

// ValidateData checks to see if the unmarshalled data is valid config.
func ValidateData(data interface{}, _ string) error {
	if data == nil {
		return fmt.Errorf("data should not be nil")
	}
	return nil
}

// ValidateConfig uses ValidateFile and ValidateData to check config for a given directory is valid Platform.sh config.
func ValidateConfig(path string) error {
	var errs error
	files := []string{".platform.app.yaml", ".platform/services.yaml", ".platform/relationships.yaml"}
	for _, file := range files {
		data, fileErr := ValidateFile(filepath.Join(path, file))
		if fileErr != nil {
			errs = multierror.Append(errs, fmt.Errorf("problem with %v: %v", file, fileErr))
		}
		dataErr := ValidateData(data, file)
		if dataErr != nil {
			errs = multierror.Append(errs, fmt.Errorf("configuration directive in %v is invalid: %v", file, dataErr))
		}
	}

	return errs
}
