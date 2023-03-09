package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/answer"
)

type DependencyManager struct{}

func (q *DependencyManager) Ask(ctx context.Context) error {
	answers, ok := answer.FromContext(ctx)
	if !ok {
		return nil
	}
	defer func() {
		ctx = answer.ToContext(ctx, answers)
	}()

	var depManagers []string
	switch answers.Stack {
	case "django":
		depManagers = []string{
			"pip",
			"poetry",
			"pipenv",
			"other",
		}
	case "laravel":
		depManagers = []string{
			"composer",
			"other",
		}
	case "next.js":
		depManagers = []string{
			"yarn",
			"npm",
			"other",
		}
	default:
		// Skip the step
		return nil
	}

	question := &survey.Select{
		Message: "Choose your dependency manager:",
		Options: depManagers,
		Default: nil,
	}

	var dependencyManager string
	err := survey.AskOne(question, &dependencyManager)
	if err != nil {
		return err
	}

	if dependencyManager != "other" {
		answers.DependencyManager = dependencyManager
	}

	return nil
}
