package platformifier

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/platformsh/platformify/internal/colors"
	"github.com/platformsh/platformify/internal/utils"
)

const (
	composerJSONFile = "composer.json"
)

func newLaravelPlatformifier(templates fs.FS, file fileCreator) *laravelPlatformifier {
	return &laravelPlatformifier{
		templates: templates,
		file:      file,
	}
}

type laravelPlatformifier struct {
	templates fs.FS
	file      fileCreator
}

func (p *laravelPlatformifier) Platformify(ctx context.Context, input *UserInput) error {
	// Get working directory.
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current working directory: %w", err)
	}

	// Check for the Laravel Bridge.
	composerJSONPaths := utils.FindAllFiles(path.Join(cwd, input.Root, input.ApplicationRoot), composerJSONFile)
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
