package question

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/models"
	"github.com/platformsh/platformify/internal/utils"
)

type DeployCommand struct{}

func (q *DeployCommand) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}
	if answers.DeployCommand != "" {
		// Skip the step
		return nil
	}

	question := &survey.Input{Message: "Deploy command:"}
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
			question.Default = fmt.Sprintf("%spython %s migrate", prefix, managePyPath)
		}
	} else if answers.Stack == models.NextJS {
		answers.DeployCommand = "./handle_mounts.sh # Move committed files from temp directory back into mounts."
		return nil
	}

	var command string
	err := survey.AskOne(question, &command)
	if err != nil {
		return err
	}

	answers.DeployCommand = command

	return nil
}
