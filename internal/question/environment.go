package question

import (
	"context"
	"fmt"

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
		var key string
		keyPrompt := survey.Input{Message: "Add an environment variable (leave blank to skip)"}
		if len(answers.Environment) > 0 {
			keyPrompt.Message = "Add another environment variable (leave blank to skip)"
		}
		if err := survey.AskOne(&keyPrompt, &key, nil); err != nil {
			return err
		}
		if key == "" {
			break
		}

		var value string
		valuePrompt := survey.Input{Message: "Set the value for \"" + key + "\""}
		if err := survey.AskOne(&valuePrompt, &value, survey.WithValidator(survey.Required)); err != nil {
			return err
		}

		answers.Environment[key] = value
	}

	return nil
}
