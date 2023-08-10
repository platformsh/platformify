package validator

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"

	"github.com/platformsh/platformify/internal/utils"
)

// ValidateFile checks the file exists and is valid yaml, then returns the unmarshalled data.
func ValidateFile(path string, schema *gojsonschema.Schema) (map[string]interface{}, error) {
	rawData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal YAML contents
	var data map[string]interface{}
	err = yaml.Unmarshal(rawData, &data)
	if err != nil {
		return nil, err
	}

	// Skip validation for empty files
	if data == nil {
		return nil, nil
	}

	result, err := schema.Validate(gojsonschema.NewGoLoader(data))
	if err != nil {
		return nil, err
	}

	if !result.Valid() {
		var validationErrors error
		for _, validationErr := range result.Errors() {
			validationErrors = errors.Join(validationErrors, errors.New(validationErr.String()))
		}
		return nil, validationErrors
	}

	return data, nil
}

// ValidateConfig uses ValidateFile and to check config for a given directory is valid Platform.sh config.
func ValidateConfig(path, flavor string) error {
	switch flavor {
	case "platformsh":
		return validatePlatformConfig(path)
	case "upsun":
		return validateUpsunConfig(path)
	}

	return fmt.Errorf("unknown flavor: %s", flavor)
}

func validatePlatformConfig(path string) error {
	var errs error
	files := map[string]*gojsonschema.Schema{
		".platform/routes.yaml":   routesSchema,
		".platform/services.yaml": servicesSchema,
	}

	for file, schema := range files {
		absPath := filepath.Join(path, file)
		if _, err := os.Stat(absPath); err != nil {
			if os.IsNotExist(err) {
				continue
			}

			errs = errors.Join(errs, fmt.Errorf("validation failed for %s: %w", file, err))
			continue
		}

		if _, err := ValidateFile(absPath, schema); err != nil {
			errs = errors.Join(errs, fmt.Errorf("validation failed for %s: %w", file, err))
		}
	}

	foundApp := false
	for _, file := range utils.FindAllFiles(path, ".platform.app.yaml") {
		foundApp = true
		if _, err := ValidateFile(file, applicationSchema); err != nil {
			relPath, _ := filepath.Rel(path, file)
			errs = errors.Join(errs, fmt.Errorf("validation failed for %s: %w", relPath, err))
		}
	}

	if errs != nil {
		return errs
	}

	if !foundApp {
		return errors.New("no application configuration found")
	}

	return nil
}

func validateUpsunConfig(path string) error {
	cnf := map[string]map[string]interface{}{}
	var errs error

	dirFs := os.DirFS(filepath.Join(path, ".upsun"))
	if err := fs.WalkDir(dirFs, ".", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() || filepath.Ext(path) != ".yaml" {
			return nil
		}

		f, err := dirFs.Open(path)
		if err != nil {
			return err
		}

		rawData, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		data := map[string]map[string]interface{}{}
		err = yaml.Unmarshal(rawData, &data)
		if err != nil {
			return err
		}

		for topKey, topValue := range data {
			if !slices.Contains([]string{"applications", "services", "routes"}, topKey) {
				errs = errors.Join(errs, fmt.Errorf("unknown key: %s", topKey))
				continue
			}

			if _, ok := cnf[topKey]; !ok {
				cnf[topKey] = topValue
				continue
			}
			for key, value := range topValue {
				if _, ok := cnf[topKey][key]; ok {
					errs = errors.Join(errs, fmt.Errorf("duplicate key: %s", key))
					continue
				}

				cnf[topKey][key] = value
			}
		}

		return nil
	}); err != nil {
		return err
	}

	result, err := upsunSchema.Validate(gojsonschema.NewGoLoader(cnf))
	if err != nil {
		errs = errors.Join(errs, err)
		return errs
	}

	if !result.Valid() {
		var validationErrors error
		for _, validationErr := range result.Errors() {
			validationErrors = errors.Join(validationErrors, errors.New(validationErr.String()))
		}
		errs = errors.Join(errs, validationErrors)
	}

	return errs
}
