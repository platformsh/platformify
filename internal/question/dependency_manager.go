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
	bundlerLockFile  = "Gemfile.lock"
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

	dependencyManagers, err := answers.Discoverer.DependencyManagers()
	if err != nil {
		return err
	}
	answers.DependencyManagers = make([]models.DepManager, 0, len(dependencyManagers))
	for _, dm := range dependencyManagers {
		answers.DependencyManagers = append(answers.DependencyManagers, models.DepManager(dm))
	}

	if exists := utils.FileExists(answers.WorkingDirectory, "", bundlerLockFile); exists {
		answers.DependencyManagers = append(answers.DependencyManagers, models.Bundler)
	}

	return nil
}
