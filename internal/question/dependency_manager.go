package question

import (
	"context"
	"os"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/models"
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
	if answers.DependencyManager.String() != "" {
		// Skip the step
		return nil
	}

	// We couldn't define the dependency manager automatically, so ask the user
	depManagers := models.DepManagersMap.Titles(answers.Type.Runtime)
	if len(depManagers) == 0 {
		// Skip the step
		return nil
	}

	question := &survey.Select{
		Message: "What dependency manager is your project using?",
		Options: depManagers,
	}

	if cwd, err := os.Getwd(); err == nil {
		switch answers.Type.Runtime {
		case models.Python:
			if exists := utils.FileExists(cwd, poetryLockFile); exists {
				question.Default = models.Poetry.Title()
			} else if exists := utils.FileExists(cwd, pipenvLockFile); exists {
				question.Default = models.Pipenv.Title()
			} else if exists := utils.FileExists(cwd, pipLockFile); exists {
				question.Default = models.Pip.Title()
			}
		case models.PHP:
			if exists := utils.FileExists(cwd, composerLockFile); exists {
				question.Default = models.Composer.Title()
			}
		case models.NodeJS:
			if exists := utils.FileExists(cwd, yarnLockFileName); exists {
				question.Default = models.Yarn.Title()
			} else if exists := utils.FileExists(cwd, npmLockFileName); exists {
				question.Default = models.Npm.Title()
			}
		default:
			// Skip the step
			return nil
		}
	}

	var title string
	if err := survey.AskOne(question, &title); err != nil {
		return err
	}

	manager, err := models.DepManagers.DepManagerByTitle(title)
	if err != nil {
		return err
	}
	answers.DependencyManager = manager

	switch manager {
	case models.Composer:
		answers.Dependencies = map[string]map[string]string{
			"php": {"composer/composer": "^2"},
		}
	case models.Npm:
		answers.Dependencies = map[string]map[string]string{
			"nodejs": {"sharp": "*"},
		}
	case models.Yarn:
		answers.Dependencies = map[string]map[string]string{
			"nodejs": {"yarn": "^1.22.0"},
		}
	}

	return nil
}
