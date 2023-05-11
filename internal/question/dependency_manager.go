package question

import (
	"context"

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
	if answers.DependencyManager.String() != "" {
		// Skip the step
		return nil
	}

	switch answers.Type.Runtime {
	case models.Python:
		if exists := utils.FileExists(answers.WorkingDirectory, poetryLockFile); exists {
			answers.DependencyManager = models.Poetry
		} else if exists := utils.FileExists(answers.WorkingDirectory, pipenvLockFile); exists {
			answers.DependencyManager = models.Pipenv
		} else if exists := utils.FileExists(answers.WorkingDirectory, pipLockFile); exists {
			answers.DependencyManager = models.Pip
		}
	case models.PHP:
		if exists := utils.FileExists(answers.WorkingDirectory, composerLockFile); exists {
			answers.DependencyManager = models.Composer
		}
	case models.NodeJS:
		if exists := utils.FileExists(answers.WorkingDirectory, yarnLockFileName); exists {
			answers.DependencyManager = models.Yarn
		} else if exists := utils.FileExists(answers.WorkingDirectory, npmLockFileName); exists {
			answers.DependencyManager = models.Npm
		}
	}

	switch answers.DependencyManager {
	case models.Composer:
		answers.Dependencies = map[string]map[string]string{
			"php": {"composer/composer": "^2"},
		}
		answers.BuildFlavor = "none"
	case models.Npm:
		answers.Dependencies = map[string]map[string]string{
			"nodejs": {"sharp": "*"},
		}
		answers.BuildFlavor = "none"
	case models.Yarn:
		answers.Dependencies = map[string]map[string]string{
			"nodejs": {"yarn": "^1.22.0"},
		}
		answers.BuildFlavor = "none"
	}

	return nil
}
