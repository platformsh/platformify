package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/answer"
)

type Type struct{}

func (q *Type) Ask(ctx context.Context) error {
	answers, ok := answer.FromContext(ctx)
	if !ok {
		return nil
	}
	defer func() {
		ctx = answer.ToContext(ctx, answers)
	}()

	if answers.Stack != "" {
		// Skip the step
		return nil
	}

	types := []string{
		"dotnet",
		"elixir",
		"golang",
		"java",
		"lisp",
		"nodejs",
		"php",
		"python",
		"ruby",
	}

	question := &survey.Select{
		Message: "Choose your PSH type:",
		Options: types,
		Default: nil,
	}

	var pshType string
	err := survey.AskOne(question, &pshType)
	if err != nil {
		return err
	}

	answers.Type = pshType

	return nil
}
