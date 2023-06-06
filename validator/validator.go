package validator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/platformsh/platformify/validator/models"
	"gopkg.in/yaml.v3"
)

// ValidateFile checks the file exists and is valid yaml, then returns the unmarshalled data.
func ValidateFile(path string) (map[string]interface{}, error) {
	// Does the file exist?
	yamlData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %v file: %v", path, err)
	}

	// Does the data unmarshal as yaml?
	var data map[string]interface{}
	err = yaml.Unmarshal(yamlData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML data for %v: %v", path, err)
	}

	// Return the yaml data.
	return data, nil
}

// ValidateData checks to see if the unmarshalled data is valid config.
func ValidateData(data map[string]interface{}, schemaFile string) error {
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
		// Validate against the model.
		_, filename := filepath.Split(schemaFile)
		// @todo update when there are models for services and routes.
		if filename == "platformsh.application.json" {
			data, readErr := os.ReadFile(filepath.Join(path, yamlFile))
			if readErr != nil {
				errs = errors.Join(errs, fmt.Errorf("could not read %v: %v", yamlFile, readErr))
			} else {
				errs = errors.Join(errs, validateAppConfig(data))
			}
		}
		// Otherwise, use the schema file.
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

func validateWithSchema(data any, schemaFile string) error {
	var errs error
	// Load and compile the JSON schema
	//schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaFile)

	//// Perform the validation
	//result, err := gojsonschema.Validate(schemaLoader, data)
	//if err != nil {
	//	errs = errors.Join(errs, fmt.Errorf("failed to validate file: %w", err))
	//}

	//if !result.Valid() {
	//	// Collect and format the validation errors
	//	for _, err := range result.Errors() {
	//		errs = errors.Join(errs, errors.New(err.String()))
	//	}
	//
	//	return fmt.Errorf("%s file is invalid:\n%s", schemaFile, errs)
	//}

	return errs
}

func validateAppConfig(data []byte) error {
	var appConfig models.AppConfig
	err := yaml.Unmarshal(data, &appConfig)
	if err != nil {
		return err
	}
	return nil
}
