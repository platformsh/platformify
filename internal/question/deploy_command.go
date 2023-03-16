package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/models"
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

	var command string
	err := survey.AskOne(question, &command)
	if err != nil {
		return err
	}

	answers.DeployCommand = command

	return nil
}
