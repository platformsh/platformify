package validator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

// ValidateFile checks the file exists and is valid yaml, then returns the unmarshalled data.
func ValidateFile(path string) ([]byte, error) {
	// Does the file exist?
	yamlData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read .platform.app.yaml file: %v", err)
	}

	// Does the yaml unmarshal properly?
	var data []byte
	err = yaml.Unmarshal(yamlData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML data: %v", err)
	}

	// Return unmarshalled data.
	return data, nil
}

// ValidateData checks to see if the unmarshalled data is valid config.
func ValidateData(data []byte, schemaFile string) error {
	if data == nil {
		return fmt.Errorf("data should not be nil")
	}
	return validateWithSchema(data, schemaFile)
}

// ValidateConfig uses ValidateFile and ValidateData to check config for a given directory is valid Platform.sh config.
func ValidateConfig(path string) error {
	var errs error
	for yamlFile, schemaFile := range schemaMap() {
		data, fileErr := ValidateFile(filepath.Join(path, yamlFile))
		if fileErr != nil {
			errs = errors.Join(errs, fmt.Errorf("problem with %v: %v", yamlFile, fileErr))
		}
		dataErr := ValidateData(data, schemaFile)
		if dataErr != nil {
			errs = errors.Join(errs, fmt.Errorf("configuration directive in %v is invalid: %v", yamlFile, dataErr))
		}
	}

	return errs
}

// The schemaMap is a hash map of yaml config file paths to their respective openAPI json schema files.
func schemaMap() map[string]string {
	schemaPath := getSchemaDir()
	var schemaMap = map[string]string{
		".platform.app.yaml":      filepath.Join(schemaPath, "platformsh.application.json"),
		".platform/routes.yaml":   filepath.Join(schemaPath, "platformsh.routes.json"),
		".platform/services.yaml": filepath.Join(schemaPath, "platformsh.services.json"),
	}
	return schemaMap
}

func getSchemaDir() string {
	dir, _ := os.Getwd()
	return filepath.Join(filepath.Dir(dir), "schema")
}

func validateWithSchema(data []byte, schemaFile string) error {
	var errs error
	// Load and compile the JSON schema
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaFile)
	documentLoader := gojsonschema.NewBytesLoader(data)

	// Perform the validation
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		errs = errors.Join(errs, fmt.Errorf("failed to validate file: %w", err))
	}

	if !result.Valid() {
		// Collect and format the validation errors
		for _, err := range result.Errors() {
			errs = errors.Join(errs, errors.New(err.String()))
		}

		return fmt.Errorf("%s file is invalid:\n%s", schemaFile, errs)
	}

	return errs
}
