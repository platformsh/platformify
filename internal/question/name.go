package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/answer"
)

type Name struct{}

func (q *Name) Ask(ctx context.Context) error {
	answers, ok := answer.FromContext(ctx)
	if !ok {
		return nil
	}
	defer func() {
		ctx = answer.ToContext(ctx, answers)
	}()

	question := &survey.Input{Message: "What is the application name?"}

	var name string
	err := survey.AskOne(question, &name)
	if err != nil {
		return err
	}

	answers.Name = name

	return nil
}
