package question

import (
	"context"
	"path"
	"path/filepath"

	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/internal/utils"
)

type ApplicationRoot struct{}

func (q *ApplicationRoot) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	for _, dm := range answers.DependencyManagers {
		switch dm {
		case models.Composer:
			if composerPath := utils.FindFile(answers.WorkingDirectory, "composer.json"); composerPath != "" {
				answers.ApplicationRoot, _ = filepath.Rel(answers.WorkingDirectory, path.Dir(composerPath))
				return nil
			}
		case models.Npm, models.Yarn:
			if packagePath := utils.FindFile(answers.WorkingDirectory, "package.json"); packagePath != "" {
				answers.ApplicationRoot, _ = filepath.Rel(answers.WorkingDirectory, path.Dir(packagePath))
				return nil
			}
		case models.Poetry:
			if pyProjectPath := utils.FindFile(answers.WorkingDirectory, "pyproject.toml"); pyProjectPath != "" {
				answers.ApplicationRoot, _ = filepath.Rel(answers.WorkingDirectory, path.Dir(pyProjectPath))
				return nil
			}
		case models.Pipenv:
			if pipfilePath := utils.FindFile(answers.WorkingDirectory, "Pipfile"); pipfilePath != "" {
				answers.ApplicationRoot, _ = filepath.Rel(answers.WorkingDirectory, path.Dir(pipfilePath))
				return nil
			}
		case models.Pip:
			if requirementsPath := utils.FindFile(answers.WorkingDirectory, "requirements.txt"); requirementsPath != "" {
				answers.ApplicationRoot, _ = filepath.Rel(answers.WorkingDirectory, path.Dir(requirementsPath))
				return nil
			}
		}
	}

	return nil
}
