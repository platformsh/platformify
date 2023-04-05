package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/models"
)

type Stack struct{}

func (q *Stack) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}
	if answers.Stack.String() != "" {
		// Skip the step
		return nil
	}

	question := &survey.Select{
		Message: "What Stack is your project using?",
		Options: models.Stacks.AllTitles(),
	}

	var title string
	err := survey.AskOne(question, &title, survey.WithPageSize(len(question.Options)))
	if err != nil {
		return err
	}
	stack, err := models.Stacks.StackByTitle(title)
	if err != nil {
		return err
	}

	answers.Stack = stack

	return nil
}
