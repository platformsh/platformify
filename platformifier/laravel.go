package platformifier

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/platformsh/platformify/internal/colors"
	"github.com/platformsh/platformify/internal/utils"
)

const (
	composerJSONFile = "composer.json"
)

func newLaravelPlatformifier(templates, fileSystem fs.FS) *laravelPlatformifier {
	return &laravelPlatformifier{
		templates:  templates,
		fileSystem: fileSystem,
	}
}

type laravelPlatformifier struct {
	templates  fs.FS
	fileSystem fs.FS
}

func (p *laravelPlatformifier) Platformify(ctx context.Context, input *UserInput) (map[string][]byte, error) {
	// Check for the Laravel Bridge.
	composerJSONPaths := utils.FindAllFiles(p.fileSystem, input.ApplicationRoot, composerJSONFile)
	for _, composerJSONPath := range composerJSONPaths {
		_, required := utils.GetJSONValue(
			p.fileSystem,
			[]string{"require", "platformsh/laravel-bridge"},
			composerJSONPath,
			true,
		)
		if !required {
			out, _, ok := colors.FromContext(ctx)
			if !ok {
				return nil, fmt.Errorf("output context failed")
			}

			var suggest = "\nPlease use composer to add the Laravel Bridge to your project:\n"
			var composerRequire = "\n    composer require platformsh/laravel-bridge\n"
			fmt.Fprintln(out, colors.Colorize(colors.WarningCode, suggest+composerRequire))
		}
	}

	return nil, nil
}
