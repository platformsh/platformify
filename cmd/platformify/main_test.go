package main

import (
	"context"
	"os"
	"testing"

	"github.com/platformsh/platformify/platformifier"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestYAMLOutput(t *testing.T) {
	// Create a temporary directory to use as the output directory.
	tempDir, err := os.MkdirTemp("", "yaml_output_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Run the command.
	ui := &platformifier.UserInput{
		ApplicationRoot: tempDir,
	}
	ctx := context.Background()
	err = runApp(ctx, ui)
	if err != nil {
		t.Fatalf("App failed with error: %v", err)
	}

	// Check the .platform.app.yaml file was created and contains valid YAML
	yamlFile := ".platform.app.yaml"
	yamlData, err := os.ReadFile(yamlFile)
	if err != nil {
		t.Fatalf("Failed to read .platform.app.yaml file: %v", err)
	}

	var data interface{}
	err = yaml.Unmarshal(yamlData, &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML data: %v", err)
	}

	// Check that the unmarshaled data is not nil
	assert.NotNil(t, data, "YAML data is nil")
}

func runApp(ctx context.Context, ui *platformifier.UserInput) error {
	return platformifier.New(ui).Platformify(ctx)
}
