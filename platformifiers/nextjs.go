package platformifiers

import (
	"context"
	"fmt"
	"os"

	"github.com/platformsh/platformify/internal/models"
	"github.com/platformsh/platformify/internal/utils"
)

const (
	nextjsTemplatesPath = "templates/nextjs"
)

type NextJSPlatformifier struct {
	*UserInput
}

func (p *NextJSPlatformifier) Platformify(ctx context.Context) error {
	if p.Stack != models.NextJS.String() {
		return fmt.Errorf("cannot platformify non-Next.js stack: %s", p.Stack)
	}

	// Gather templates.
	templates, err := utils.GatherTemplates(ctx, templatesFs, nextjsTemplatesPath)
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
