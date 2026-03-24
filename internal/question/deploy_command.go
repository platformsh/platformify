package question

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"

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
		managePyPath := utils.FindFile(answers.WorkingDirectory, answers.ApplicationRoot, managePyFile)
		if managePyPath != "" {
			managePyPath, _ = filepath.Rel(answers.ApplicationRoot, managePyPath)
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
			"mkdir -p storage/framework/sessions",
			"mkdir -p storage/framework/cache",
			"mkdir -p storage/framework/views",
			"php artisan migrate --force",
			"php artisan optimize:clear",
		)
	case models.Rails:
		answers.DeployCommand = append(answers.DeployCommand,
			"bundle exec rake db:migrate",
		)
	}

	return nil
}
