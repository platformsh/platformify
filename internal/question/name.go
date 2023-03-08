package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"
)

type Name Question

func (q *Name) Ask(ctx context.Context) error {
	question := &survey.Input{Message: "What is the application name?"}

	var name string
	err := survey.AskOne(question, &name)
	if err != nil {
		return err
	}

	q.Answers.Name = name

	return nil
}
