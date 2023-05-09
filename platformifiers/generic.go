package platformifiers

import (
	"context"
	"fmt"

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

	if err := utils.WriteTemplates(ctx, p.WorkingDirectory, templates, p.UserInput); err != nil {
		return fmt.Errorf("could not write Platform.sh files: %w", err)
	}

	return nil
}
