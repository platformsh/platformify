package question

import (
	"context"
	"fmt"

	"github.com/platformsh/platformify/internal/colors"
	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/internal/utils"
)

const (
	npmLockFileName  = "package-lock.json"
	yarnLockFileName = "yarn.lock"
	poetryLockFile   = "poetry.lock"
	pipenvLockFile   = "Pipfile.lock"
	pipLockFile      = "requirements.txt"
	composerLockFile = "composer.lock"
)

type DependencyManager struct{}

func (q *DependencyManager) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}
	if len(answers.DependencyManagers) > 0 {
		// Skip the step
		return nil
	}

	answers.Dependencies = map[string]map[string]string{}
	answers.BuildFlavor = "none"

	defer func() {
		_, stderr, ok := colors.FromContext(ctx)
		if !ok {
			return
		}

		if len(answers.DependencyManagers) > 0 {
			dependencyManagers := answers.DependencyManagers[0].Title()
			for _, dm := range answers.DependencyManagers[1:] {
				dependencyManagers = fmt.Sprintf("%s, %s", dependencyManagers, dm.Title())
			}

			fmt.Fprintf(
				stderr,
				"%s %s\n",
				colors.Colorize(colors.GreenCode, "✓"),
				colors.Colorize(
					colors.BrandCode,
					fmt.Sprintf("Detected dependency managers: %s", dependencyManagers),
				),
			)
		}
	}()

	if exists := utils.FileExists(answers.WorkingDirectory, poetryLockFile); exists {
		answers.DependencyManagers = append(answers.DependencyManagers, models.Poetry)
	} else if exists := utils.FileExists(answers.WorkingDirectory, pipenvLockFile); exists {
		answers.DependencyManagers = append(answers.DependencyManagers, models.Pipenv)
	} else if exists := utils.FileExists(answers.WorkingDirectory, pipLockFile); exists {
		answers.DependencyManagers = append(answers.DependencyManagers, models.Pip)
	}

	if exists := utils.FileExists(answers.WorkingDirectory, composerLockFile); exists {
		answers.DependencyManagers = append(answers.DependencyManagers, models.Composer)
		answers.Dependencies["php"] = map[string]string{"composer/composer": "^2"}
	}

	if exists := utils.FileExists(answers.WorkingDirectory, yarnLockFileName); exists {
		answers.DependencyManagers = append(answers.DependencyManagers, models.Yarn)
		answers.Dependencies["nodejs"] = map[string]string{"yarn": "^1.22.0"}
	} else if exists := utils.FileExists(answers.WorkingDirectory, npmLockFileName); exists {
		answers.DependencyManagers = append(answers.DependencyManagers, models.Npm)
	}

	return nil
}
