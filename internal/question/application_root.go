package question

import (
	"context"
	"os"
	"path"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/models"
	"github.com/platformsh/platformify/internal/utils"
)

type ApplicationRoot struct{}

func (q *ApplicationRoot) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}
	if answers.ApplicationRoot != "" {
		// Skip the step
		return nil
	}

	question := &survey.Input{Message: "Application root:", Default: "."}
	cwd, _ := os.Getwd()
	switch answers.DependencyManager {
	case models.Composer:
		if composerPath := utils.FindFile(cwd, "composer.json"); composerPath != "" {
			question.Default, _ = filepath.Rel(cwd, path.Dir(composerPath))
		}
	case models.Npm, models.Yarn:
		if packagePath := utils.FindFile(cwd, "package.json"); packagePath != "" {
			question.Default, _ = filepath.Rel(cwd, path.Dir(packagePath))
		}
	case models.Poetry:
		if pyProjectPath := utils.FindFile(cwd, "pyproject.toml"); pyProjectPath != "" {
			question.Default, _ = filepath.Rel(cwd, path.Dir(pyProjectPath))
		}
	case models.Pipenv:
		if pipfilePath := utils.FindFile(cwd, "Pipfile"); pipfilePath != "" {
			question.Default, _ = filepath.Rel(cwd, path.Dir(pipfilePath))
		}
	case models.Pip:
		if requirementsPath := utils.FindFile(cwd, "requirements.txt"); requirementsPath != "" {
			question.Default, _ = filepath.Rel(cwd, path.Dir(requirementsPath))
		}
	}

	return nil
}
