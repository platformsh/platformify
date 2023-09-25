package question

import (
	"context"
	"fmt"
	"path"
	"path/filepath"

	"golang.org/x/exp/slices"

	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/internal/utils"
)

type DeployCommand struct{}

func (q *DeployCommand) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	switch answers.Stack {
	case models.Django:
		managePyPath := utils.FindFile(path.Join(answers.WorkingDirectory, answers.ApplicationRoot), managePyFile)
		if managePyPath != "" {
			managePyPath, _ = filepath.Rel(path.Join(answers.WorkingDirectory, answers.ApplicationRoot), managePyPath)
			prefix := ""
			if slices.Contains(answers.DependencyManagers, models.Pipenv) {
				prefix = "pipenv run "
			} else if slices.Contains(answers.DependencyManagers, models.Poetry) {
				prefix = "poetry run "
			}
			answers.DeployCommand = append(answers.DeployCommand,
				fmt.Sprintf("%spython %s migrate", prefix, managePyPath),
			)
		}
	case models.Laravel:
		answers.DeployCommand = append(answers.DeployCommand,
			"mkdir -p storage/framework/{sessions,cache,views}",
			"php artisan optimize:clear",
			"php artisan migrate --force",
		)
	}

	return nil
}
