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
	for _, stack := range models.Stacks {
		// Create a temporary directory to use as the output directory.
		tempDir, err := os.MkdirTemp("", stack.Title())
		if err != nil {
			t.Fatalf("Failed to create temporary directory: %v", err)
		}

		// Run the command.
		ui := &platformifier.UserInput{
			Stack:            platformifier.Stack(stack),
			WorkingDirectory: tempDir,
		}
		ctx := context.Background()
		err = runApp(ctx, ui)
		if err != nil {
			t.Fatalf("App failed with error: %v", err)
		}

		// Validate the config.
		invalid := validator.ValidateConfig(tempDir)
		assert.NoError(t, err, "error while validating config: %v", invalid)

		rmErr := os.RemoveAll(tempDir)
		if rmErr != nil {
			t.Fatalf("Could not remove temporary directory: %v", tempDir)
		}
	}
}

func runApp(ctx context.Context, ui *platformifier.UserInput) error {
	return platformifier.New(ui).Platformify(ctx)
}
