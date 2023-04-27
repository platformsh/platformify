package platformifiers

import (
	"context"
	"fmt"
	"os"

	"github.com/platformsh/platformify/internal/models"
	"github.com/platformsh/platformify/internal/utils"
)

const laravelTemplatesPath = "templates/laravel"

type LaravelPlatformifier struct {
	*UserInput
}

func (p *LaravelPlatformifier) Platformify(ctx context.Context) error {
	if p.Stack != models.Laravel.String() {
		return fmt.Errorf("cannot platformify non-Laravel stack: %s", p.Stack)
	}

	// Gather templates.
	templates, err := utils.GatherTemplates(ctx, templatesFs, laravelTemplatesPath)
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
