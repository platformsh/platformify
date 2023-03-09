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

	if answers.Type.Name != "" {
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

	var name string
	err := survey.AskOne(question, &name, survey.WithPageSize(len(question.Options)))
	if err != nil {
		return err
	}

	switch name {
	case "dotnet":
		question = &survey.Select{
			Message: "Choose C#/.Net Core version:",
			Options: []string{
				"6.0", "3.1",
			},
			Default: "6.0",
		}
	case "elixir":
		question = &survey.Select{
			Message: "Choose Elixir version:",
			Options: []string{
				"1.13", "1.12", "1.11",
			},
			Default: "1.13",
		}
	case "golang":
		question = &survey.Select{
			Message: "Choose Go version:",
			Options: []string{
				"1.20", "1.19",
			},
			Default: "1.20",
		}
	case "java":
		question = &survey.Select{
			Message: "Choose Java version:",
			Options: []string{
				"19", "18", "17", "11", "8",
			},
			Default: "19",
		}
	case "lisp":
		question = &survey.Select{
			Message: "Choose Lisp version:",
			Options: []string{
				"2.1", "2.0", "1.5",
			},
			Default: "2.1",
		}
	case "nodejs":
		question = &survey.Select{
			Message: "Choose JavaScript/Node.js version:",
			Options: []string{
				"18", "16", "14",
			},
			Default: "18",
		}
	case "php":
		question = &survey.Select{
			Message: "Choose PHP version:",
			Options: []string{
				"8.2", "8.1", "8.0",
			},
			Default: "8.2",
		}
	case "python":
		question = &survey.Select{
			Message: "Choose Python version:",
			Options: []string{
				"3.11", "3.10", "3.9", "3.8", "3.7",
			},
			Default: "3.11",
		}
	case "ruby":
		question = &survey.Select{
			Message: "Choose Ruby version:",
			Options: []string{
				"3.2", "3.1", "3.0", "2.7", "2.6", "2.5", "2.4", "2.3",
			},
			Default: "3.2",
		}
	default:
		return nil
	}

	var version string
	err = survey.AskOne(question, &version, survey.WithPageSize(len(question.Options)))
	if err != nil {
		return err
	}

	answers.Type.Name = name
	answers.Type.Version = version

	return nil
}
