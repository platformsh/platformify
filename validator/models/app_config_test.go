package models

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

// Test yaml export of an AppConfig object.
func TestAppConfigExport(t *testing.T) {
	// Create an AppConfig Object.
	appConfig := getAppConfig()

	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "app_config_yaml_tests")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Export the AppConfig yaml into the temporary directory.
	yamlData, err := yaml.Marshal(appConfig)
	assert.NoError(t, err, "failed to marshal AppConfig YAML: %v", err)

	filePath := path.Join(tempDir, "app-config.yaml")
	err = os.WriteFile(filePath, yamlData, 0o644)
	assert.NoError(t, err, "failed to write YAML file: %v", err)

	// Load the yaml from the file.
	loaded, loadErr := os.ReadFile(filePath)
	assert.NoError(t, loadErr, "failed to load app config yaml file")
	var gotAppConfig AppConfig
	unmarshalErr := yaml.Unmarshal(loaded, &gotAppConfig)
	assert.NoError(t, unmarshalErr, "failed to unmarshal yaml into AppConfig object")

	// Compare to the original object.
	assert.Equalf(t, appConfig, &gotAppConfig, "AppConfig objects are not equal")
}

func getAppConfig() *AppConfig {
	appType, _ := NewAppType("golang", "1.20")
	// Create an AppConfig object with ALL THE THINGS!
	return &AppConfig{
		Name: "TestAppConfig",
		Type: appType,
		Size: "XL",
		Disk: 128,
		//Build:           nil,
		//Relationships:   nil,
		//Variables:       nil,
		//Hooks:           nil,
		//Dependencies:    nil,
		//Mounts:          nil,
		//Web:             nil,
		//Crons:           nil,
		//Workers:         nil,
		//Timezone:        "",
		//Access:          nil,
		//Firewall:        nil,
		//Source:          nil,
		//Runtime:         nil,
		//AdditionalHosts: nil,
	}
}
