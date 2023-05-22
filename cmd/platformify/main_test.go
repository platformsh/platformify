package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/platformsh/platformify/internal/validator"
	"github.com/platformsh/platformify/platformifier"

	"github.com/stretchr/testify/assert"
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
		WorkingDirectory: tempDir,
	}
	ctx := context.Background()
	err = runApp(ctx, ui)
	if err != nil {
		t.Fatalf("App failed with error: %v", err)
	}

	// Check the .platform.app.yaml file was created and contains valid YAML
	yamlFile := filepath.Join(tempDir, ".platform.app.yaml")
	yamlData, err := validator.ValidateFile(yamlFile)
	assert.NoError(t, err, "error while validating .platform.app.yaml")
	assert.NoError(t, validator.ValidateData(yamlData))
}

func runApp(ctx context.Context, ui *platformifier.UserInput) error {
	return platformifier.New(ui).Platformify(ctx)
}
