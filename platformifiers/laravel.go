package platformifiers

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/platformsh/platformify/internal/colors"
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

	// Check for the Laravel Bridge.
	appRoot := path.Join(cwd, p.Root, p.ApplicationRoot)
	composerJSONPaths := utils.FindAllFiles(appRoot, "composer.json")
	for _, composerJSONPath := range composerJSONPaths {
		_, required := utils.GetJSONKey([]string{"require", "platformsh/laravel-bridge"}, composerJSONPath)
		if !required {
			out, _, ok := colors.FromContext(ctx)
			if !ok {
				return fmt.Errorf("output context failed")
			}

			var suggest = "\nPlease use composer to add the Laravel Bridge to your project:\n"
			var composerRequire = "\n    composer require platformsh/laravel-bridge\n\n"
			fmt.Fprint(out, colors.Colorize(colors.WarningCode, suggest+composerRequire))
		}
	}

	return nil
}
