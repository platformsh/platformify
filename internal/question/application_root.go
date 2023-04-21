package question

import (
	"context"
	"os"
	"path"
	"path/filepath"

	"github.com/platformsh/platformify/internal/models"
	"github.com/platformsh/platformify/internal/utils"
)

type ApplicationRoot struct{}

func (q *ApplicationRoot) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}
	cwd, _ := os.Getwd()
	switch answers.DependencyManager {
	case models.Composer:
		if composerPath := utils.FindFile(cwd, "composer.json"); composerPath != "" {
			answers.ApplicationRoot, _ = filepath.Rel(cwd, path.Dir(composerPath))
		}
	case models.Npm, models.Yarn:
		if packagePath := utils.FindFile(cwd, "package.json"); packagePath != "" {
			answers.ApplicationRoot, _ = filepath.Rel(cwd, path.Dir(packagePath))
		}
	case models.Poetry:
		if pyProjectPath := utils.FindFile(cwd, "pyproject.toml"); pyProjectPath != "" {
			answers.ApplicationRoot, _ = filepath.Rel(cwd, path.Dir(pyProjectPath))
		}
	case models.Pipenv:
		if pipfilePath := utils.FindFile(cwd, "Pipfile"); pipfilePath != "" {
			answers.ApplicationRoot, _ = filepath.Rel(cwd, path.Dir(pipfilePath))
		}
	case models.Pip:
		if requirementsPath := utils.FindFile(cwd, "requirements.txt"); requirementsPath != "" {
			answers.ApplicationRoot, _ = filepath.Rel(cwd, path.Dir(requirementsPath))
		}
	}

	return nil
}
