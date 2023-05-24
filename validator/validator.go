package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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
func ValidateData(data interface{}, schemaFile string) error {
	if data == nil {
		return fmt.Errorf("data should not be nil")
	}

}

// ValidateConfig uses ValidateFile and ValidateData to check config for a given directory is valid Platform.sh config.
func ValidateConfig(path string) error {
	var errs error
	for yamlFile, schemaFile := range schemaMap() {
		data, fileErr := ValidateFile(filepath.Join(path, yamlFile))
		if fileErr != nil {
			errs = multierror.Append(errs, fmt.Errorf("problem with %v: %v", yamlFile, fileErr))
		}
		dataErr := ValidateData(data, schemaFile)
		if dataErr != nil {
			errs = multierror.Append(errs, fmt.Errorf("configuration directive in %v is invalid: %v", yamlFile, dataErr))
		}
	}

	return errs
}

// The schemaMap is a hash map of yaml config file paths to their respective openAPI json schema files.
func schemaMap() map[string]string {
	schemaPath := getSchemaDir()
	var schemaMap = map[string]string{
		".platform.app.yaml":      filepath.Join(schemaPath, "platformsh.application.json"),
		".platform/services.yaml": filepath.Join(schemaPath, "platformsh.services.json"),
		".platform/routes.yaml":   filepath.Join(schemaPath, "platformsh.routes.json"),
	}
	return schemaMap
}

func getSchemaDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Join(filepath.Dir(filename), "schema")
}
