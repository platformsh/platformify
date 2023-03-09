package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/answer"
)

type Stack struct{}

func (q *Stack) Ask(ctx context.Context) error {
	answers, ok := answer.FromContext(ctx)
	if !ok {
		return nil
	}
	defer func() {
		ctx = answer.ToContext(ctx, answers)
	}()

	stacks := []string{
		"django",
		"laravel",
		"next.js",
		"other",
	}

	question := &survey.Select{
		Message: "Choose your stack:",
		Options: stacks,
		Default: nil,
	}

	var stack string
	err := survey.AskOne(question, &stack, survey.WithPageSize(len(question.Options)))
	if err != nil {
		return err
	}

	if stack == "other" {
		answers.Stack = "generic"
	} else {
		var name string
		switch stack {
		case "django":
			name = "python"
		case "laravel":
			name = "php"
		case "next.js":
			name = "nodejs"
		}

		answers.Stack = stack
		answers.Type.Name = name
	}

	return nil
}
