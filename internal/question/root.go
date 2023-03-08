package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"
)

type Root Question

func (q *Root) Ask(ctx context.Context) error {
	question := &survey.Input{Message: "What is the application root?"}

	var root string
	err := survey.AskOne(question, &root)
	if err != nil {
		return err
	}

	q.Answers.Root = root

	return nil
}
