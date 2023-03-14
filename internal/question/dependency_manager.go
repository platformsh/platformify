package question

import (
	"context"
	"os"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/models"
)

const (
	npmLockFileName  = "package-lock.json"
	yarnLockFileName = "yarn.lock"
)

type DependencyManager struct{}

func (q *DependencyManager) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	switch answers.Type.Runtime {
	case models.Python:
		// TODO: check for python dependency manager
	case models.PHP:
		// TODO: check for php dependency manager
	case models.NodeJS:
		// Check if the project uses "npm" as a dependency manager
		exists, err := fileExists(npmLockFileName)
		if err != nil {
			return err
		} else if exists {
			answers.DependencyManager = models.Npm
			return nil
		}

		// Check if the project uses "yarn" as a dependency manager
		exists, err = fileExists(yarnLockFileName)
		if err != nil {
			return err
		} else if exists {
			answers.DependencyManager = models.Yarn
			return nil
		}
	default:
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
		Message: "Choose your dependency manager:",
		Options: depManagers,
	}

	var title string
	err := survey.AskOne(question, &title)
	if err != nil {
		return err
	}

	manager, err := models.DepManagers.DepManagerByTitle(title)
	if err != nil {
		return err
	}
	answers.DependencyManager = manager

	return nil
}

// fileExists checks if the file exists
func fileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
