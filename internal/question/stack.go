package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/question/models"
)

type Stack struct{}

func (q *Stack) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	question := &survey.Select{
		Message: "What Stack is your project using?",
		Options: models.Stacks.AllTitles(),
	}

	var stack models.Stack
	err := survey.AskOne(question, &stack, survey.WithPageSize(len(question.Options)))
	if err != nil {
		return err
	}

	answers.Stack = stack

	return nil
}
