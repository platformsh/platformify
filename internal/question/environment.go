package question

import (
	"context"

	"github.com/platformsh/platformify/internal/models"
)

type Environment struct{}

func (q *Environment) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
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

	return nil
}
