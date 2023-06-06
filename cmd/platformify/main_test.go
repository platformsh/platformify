package main

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/platformifier"
	"github.com/platformsh/platformify/validator"
)

func TestYAMLOutput(t *testing.T) {
	// Create a temporary directory to use as the output directory.
	tempDir, err := os.MkdirTemp("", "yaml_tests")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	for _, stack := range models.Stacks {
		dir, err1 := os.MkdirTemp(tempDir, stack.Title())
		if err1 != nil {
			t.Fatalf("Failed to create temporary %v directory: %v", stack.Title(), err1)
		}
		// Run the command.
		ui := &platformifier.UserInput{
			Name:             stack.Title() + "Test",
			Type:             "python: \"3.11\"",
			Stack:            platformifier.Stack(stack),
			WorkingDirectory: dir,
		}
		ctx := context.Background()
		err = runApp(ctx, ui)
		if err != nil {
			t.Fatalf("App failed with error: %v", err)
		}

		// Validate the config.
		invalid := validator.ValidateConfig(dir)
		assert.NoError(t, invalid, "error while validating config: %v", invalid)
	}
}

func runApp(ctx context.Context, ui *platformifier.UserInput) error {
	return platformifier.New(ui).Platformify(ctx)
}
