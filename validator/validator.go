package validator

import (
	"fmt"
	"os"

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
func ValidateData(data interface{}) error {
	if data == nil {
		return fmt.Errorf("data should not be nil")
	}
	return nil
}
