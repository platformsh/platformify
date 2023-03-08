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
		"Django",
		"Laravel",
		"Next.js",
		"Other",
	}

	question := &survey.Select{
		Message: "Choose your stack:",
		Options: stacks,
		Default: nil,
	}

	var stack string
	err := survey.AskOne(question, &stack)
	if err != nil {
		return err
	}

	if stack != "Other" {
		var pshType string
		switch stack {
		case "Django":
			pshType = "python"
		case "Laravel":
			pshType = "php"
		case "Next.js":
			pshType = "nodejs"
		}

		answers.Stack = stack
		answers.Type = pshType
	}

	return nil
}
