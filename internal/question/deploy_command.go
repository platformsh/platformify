package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"
)

type DeployCommand Question

func (q *DeployCommand) Ask(ctx context.Context) error {
	question := &survey.Input{Message: "Deploy command?"}

	var deployCommand string
	err := survey.AskOne(question, &deployCommand)
	if err != nil {
		return err
	}

	q.Answers.DeployCommand = deployCommand

	return nil
}
