package question

import (
	"context"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/models"
)

type Environment struct{}

func (q *Environment) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}
	if len(answers.Environment) != 0 {
		// Skip the step
		return nil
	}

	if answers.Stack == models.Django {
		if answers.DependencyManager == models.Poetry {
			answers.Environment = map[string]string{
				"POETRY_VERSION":                "1.4.0",
				"POETRY_VIRTUALENVS_IN_PROJECT": "true",
			}
		} else if answers.DependencyManager == models.Pipenv {
			answers.Environment = map[string]string{
				"PIPENV_VERSION":         "2023.2.18",
				"PIPENV_VENV_IN_PROJECT": "1",
			}
		}
	}

	if len(answers.Environment) > 0 {
		fmt.Println("We identified a few environment variables for you already!")
		for key, value := range answers.Environment {
			fmt.Println("  ", key, "=", value)
		}
	}

	for {
		question := &survey.Confirm{
			Message: "Do you want to add a variable?",
			Default: true,
		}

		var addVariable = true
		err := survey.AskOne(question, &addVariable)
		if err != nil {
			return err
		}

		if !addVariable {
			break
		}

		questions := []*survey.Question{
			{
				Name:      "key",
				Prompt:    &survey.Input{Message: "Env var name"},
				Validate:  survey.Required,
				Transform: survey.TransformString(strings.ToUpper),
			},
			{
				Name:   "value",
				Prompt: &survey.Input{Message: "Env var value"},
			},
		}

		variable := struct {
			Key   string
			Value string
		}{}

		err = survey.Ask(questions, &variable)
		if err != nil {
			return err
		}

		answers.Environment[variable.Key] = variable.Value
	}

	return nil
}
