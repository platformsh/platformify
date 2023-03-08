package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"
)

type DependencyManager Question

func (q *DependencyManager) Ask(ctx context.Context) error {
	var depManagers []string
	switch q.Answers.Stack {
	case "Django":
		depManagers = []string{
			"pip",
			"poetry",
			"pipenv",
			"other",
		}
	case "Laravel":
		depManagers = []string{
			"composer",
			"other",
		}
	case "Next.js":
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
		q.Answers.DependencyManager = dependencyManager
	}

	return nil
}
