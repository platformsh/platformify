package validator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/platformsh/platformify/internal/utils"
)

// ValidateFile checks the file exists and is valid yaml, then returns the unmarshalled data.
func ValidateFile(path string) (map[string]interface{}, error) {
	// Does the file exist?
	rawData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Does the yamlData unmarshal as yaml?
	var yamlData map[string]interface{}
	err = yaml.Unmarshal(rawData, &yamlData)
	if err != nil {
		return nil, err
	}

	// Return the unmarshalled yaml yamlData.
	return yamlData, nil
}

// ValidateConfig uses ValidateFile and to check config for a given directory is valid Platform.sh config.
func ValidateConfig(path string) error {
	var errs error
	files := [3]string{".platform/applications.yaml", ".platform/routes.yaml", ".platform/services.yaml"}
	for _, file := range files {
		absPath := filepath.Join(path, file)
		if _, err := os.Stat(absPath); err != nil {
			if os.IsNotExist(err) {
				continue
			}

			errs = errors.Join(errs, fmt.Errorf("validation failed for %s: %w", file, err))
			continue
		}

		if _, err := ValidateFile(absPath); err != nil {
			errs = errors.Join(errs, fmt.Errorf("validation failed for %s: %w", file, err))
		}
	}

	for _, file := range utils.FindAllFiles(path, ".platform.app.yaml") {
		fmt.Println(path)
		if _, err := ValidateFile(file); err != nil {
			relPath, _ := filepath.Rel(path, file)
			errs = errors.Join(errs, fmt.Errorf("validation failed for %s: %w", relPath, err))
		}
	}

	return errs
}
