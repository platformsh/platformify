package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"
)

type Stack Question

func (q *Stack) Ask(ctx context.Context) error {
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

		q.Answers.Stack = stack
		q.Answers.Type = pshType
	}

	return nil
}
