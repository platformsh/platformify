package platformifiers

import (
	"context"
	"fmt"
	"os"

	"github.com/platformsh/platformify/internal/utils"
)

const genericTemplatesPath = "templates/generic"

// GenericPlatformifier contains the configuration for the application to Platformify
type GenericPlatformifier struct {
	*UserInput
}

// Platformify will generate the .platform.app.yaml and .platform/ directory
func (p *GenericPlatformifier) Platformify(ctx context.Context) error {
	// Gather templates.
	templates, err := utils.GatherTemplates(ctx, templatesFs, genericTemplatesPath)
	if err != nil {
		return err
	}

	// Get working directory.
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current working directory: %w", err)
	}
	if err := utils.WriteTemplates(ctx, cwd, templates, p.UserInput); err != nil {
		return fmt.Errorf("could not write Platform.sh files: %w", err)
	}

	return nil
}
