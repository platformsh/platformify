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

	deployCommand := ""
	cwd, _ := os.Getwd()
	if answers.Stack == models.Django {
		if managePyPath := utils.FindFile(path.Join(cwd, answers.ApplicationRoot), managePyFile); managePyPath != "" {
			managePyPath, _ = filepath.Rel(path.Join(cwd, answers.ApplicationRoot), managePyPath)
			prefix := ""
			switch answers.DependencyManager {
			case models.Pipenv:
				prefix = "pipenv run "
			case models.Poetry:
				prefix = "poetry run "
			}
			deployCommand = fmt.Sprintf("%spython %s migrate", prefix, managePyPath)
		}
	}

	if deployCommand != "" {
		fmt.Println("We identified the command to use during deployment for you!")
		fmt.Println("  ", deployCommand)
	}

	return nil
}
