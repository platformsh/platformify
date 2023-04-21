package question

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/internal/utils"
)

type DeployCommand struct{}

func (q *DeployCommand) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	cwd, _ := os.Getwd()
	switch answers.Stack {
	case models.Django:
		if managePyPath := utils.FindFile(path.Join(cwd, answers.ApplicationRoot), managePyFile); managePyPath != "" {
			managePyPath, _ = filepath.Rel(path.Join(cwd, answers.ApplicationRoot), managePyPath)
			prefix := ""
			switch answers.DependencyManager {
			case models.Pipenv:
				prefix = "pipenv run "
			case models.Poetry:
				prefix = "poetry run "
			}
			answers.DeployCommand = append(answers.DeployCommand,
				fmt.Sprintf("%spython %s migrate", prefix, managePyPath),
			)
		}
	case models.Laravel:
		answers.DeployCommand = append(answers.DeployCommand,
			"php artisan optimize:clear",
			"php artisan migrate --force",
		)
	}

	return nil
}
