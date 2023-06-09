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
	if len(answers.DependencyManagers) > 0 {
		// Skip the step
		return nil
	}

	answers.Dependencies = map[string]map[string]string{}
	answers.BuildFlavor = "none"

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
		answers.Dependencies["nodejs"] = map[string]string{"sharp": "*"}
	}

	return nil
}
