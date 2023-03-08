package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/answer"
)

type DeployCommand struct{}

func (q *DeployCommand) Ask(ctx context.Context) error {
	answers, ok := answer.FromContext(ctx)
	if !ok {
		return nil
	}
	defer func() {
		ctx = answer.ToContext(ctx, answers)
	}()

	question := &survey.Input{Message: "Deploy command?"}

	var deployCommand string
	err := survey.AskOne(question, &deployCommand)
	if err != nil {
		return err
	}

	answers.DeployCommand = deployCommand

	return nil
}
