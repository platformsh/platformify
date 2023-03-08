package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"
)

type Listen Question

func (q *Listen) Ask(ctx context.Context) error {
	interfaces := []string{
		"HTTP",
		"Unix-socket",
	}

	question := &survey.Select{
		Message: "Choose interface to listen to:",
		Options: interfaces,
		Default: nil,
	}

	var listen string
	err := survey.AskOne(question, &listen)
	if err != nil {
		return err
	}

	q.Answers.Listen = listen

	return nil
}
