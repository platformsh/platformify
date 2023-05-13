package platformifier

import (
	"context"
	"fmt"
	"io/fs"
	"path"

	"github.com/platformsh/platformify/internal/colors"
	"github.com/platformsh/platformify/internal/utils"
)

const (
	composerJSONFile = "composer.json"
)

func newLaravelPlatformifier(templates fs.FS) *laravelPlatformifier {
	return &laravelPlatformifier{
		templates: templates,
	}
}

type laravelPlatformifier struct {
	templates fs.FS
}

func (p *laravelPlatformifier) Platformify(ctx context.Context, input *UserInput) error {
	// Check for the Laravel Bridge.
	appRoot := path.Join(input.WorkingDirectory, input.Root, input.ApplicationRoot)
	composerJSONPaths := utils.FindAllFiles(appRoot, composerJSONFile)
	for _, composerJSONPath := range composerJSONPaths {
		_, required := utils.GetJSONKey([]string{"require", "platformsh/laravel-bridge"}, composerJSONPath)
		if !required {
			out, _, ok := colors.FromContext(ctx)
			if !ok {
				return fmt.Errorf("output context failed")
			}

			var suggest = "\nPlease use composer to add the Laravel Bridge to your project:\n"
			var composerRequire = "\n    composer require platformsh/laravel-bridge\n"
			fmt.Fprintln(out, colors.Colorize(colors.WarningCode, suggest+composerRequire))
		}
	}

	return nil
}
