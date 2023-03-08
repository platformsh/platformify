package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"
)

type Type Question

func (q *Type) Ask(ctx context.Context) error {
	if q.Answers.Stack != "" {
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

	q.Answers.Type = pshType

	return nil
}
