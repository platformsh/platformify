package validator

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"

	"github.com/platformsh/platformify/internal/utils"
)

// ValidateFile checks the file exists and is valid yaml, then returns the unmarshalled data.
func ValidateFile(fileSystem fs.FS, path string, schema *gojsonschema.Schema) (map[string]interface{}, error) {
	rawData, err := fs.ReadFile(fileSystem, path)
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

// ValidateConfig uses ValidateFile and to check config for a given directory is valid config.
func ValidateConfig(fileSystem fs.FS, flavor string) error {
	switch flavor {
	case "platform":
		return validatePlatformConfig(fileSystem)
	case "upsun":
		return validateUpsunConfig(fileSystem)
	default:
		return fmt.Errorf("unknown flavor: %s", flavor)
	}
}

func validatePlatformConfig(fileSystem fs.FS) error {
	var errs error
	files := map[string]*gojsonschema.Schema{
		".platform/routes.yaml":   routesSchema,
		".platform/services.yaml": servicesSchema,
	}

	for file, schema := range files {
		if _, err := fs.Stat(fileSystem, file); err != nil {
			if os.IsNotExist(err) {
				continue
			}

			errs = errors.Join(errs, fmt.Errorf("validation failed for %s: %w", file, err))
			continue
		}

		if _, err := ValidateFile(fileSystem, file, schema); err != nil {
			errs = errors.Join(errs, fmt.Errorf("validation failed for %s: %w", file, err))
		}
	}

	foundApp := false
	for _, file := range utils.FindAllFiles(fileSystem, "", ".platform.app.yaml") {
		foundApp = true
		if _, err := ValidateFile(fileSystem, file, applicationSchema); err != nil {
			errs = errors.Join(errs, fmt.Errorf("validation failed for %s: %w", file, err))
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

func validateUpsunConfig(dirFs fs.FS) error {
	cnf := map[string]map[string]interface{}{}
	var errs error

	if stat, err := fs.Stat(dirFs, ".upsun"); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("the .upsun directory does not exist")
		}
		return fmt.Errorf("cannot open the .upsun directory")
	} else if !stat.IsDir() {
		return fmt.Errorf(".upsun is not a directory")
	}
	entries, err := fs.ReadDir(dirFs, ".upsun")
	if err != nil {
		return fmt.Errorf("cannot open the .upsun directory")
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}

		f, loopErr := dirFs.Open(filepath.Join(".upsun", entry.Name()))
		if loopErr != nil {
			return loopErr
		}

		rawData, loopErr := io.ReadAll(f)
		if loopErr != nil {
			return loopErr
		}

		data := map[string]map[string]interface{}{}
		loopErr = yaml.Unmarshal(rawData, &data)
		if loopErr != nil {
			return fmt.Errorf("unmarshal failed for %s: %w", entry.Name(), loopErr)
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
