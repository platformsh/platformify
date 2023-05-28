package question

import (
	"context"
	"fmt"
	"path"
	"path/filepath"

	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/internal/utils"
)

type BuildSteps struct{}

func (q *BuildSteps) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	switch answers.DependencyManager {
	case models.Poetry:
		answers.BuildSteps = append(
			answers.BuildSteps,
			"# Set PIP_USER to 0 so that Poetry does not complain",
			"export PIP_USER=0",
			"# Install poetry as a global tool",
			"python -m venv /app/.global",
			"pip install poetry==$POETRY_VERSION",
			"poetry install",
		)
	case models.Pipenv:
		answers.BuildSteps = append(
			answers.BuildSteps,
			"# Set PIP_USER to 0 so that Pipenv does not complain",
			"export PIP_USER=0",
			"# Install Pipenv as a global tool",
			"python -m venv /app/.global",
			"pip install pipenv==$PIPENV_VERSION",
			"pipenv install",
		)
	case models.Pip:
		answers.BuildSteps = append(
			answers.BuildSteps,
			"pip install -r requirements.txt",
		)

	case models.Yarn:
		answers.BuildSteps = append(
			answers.BuildSteps,
			"yarn",
			"yarn build",
		)
	case models.Npm:
		answers.BuildSteps = append(
			answers.BuildSteps,
			"npm i",
			"npm run build",
		)
	case models.Composer:
		answers.BuildSteps = append(
			answers.BuildSteps,
			"composer --no-ansi --no-interaction install --no-progress --prefer-dist --optimize-autoloader --no-dev",
			"# Install a specific NodeJS version https://github.com/platformsh/snippets/",
			"# uncomment next line to build assets deploying",
			"# npm install && npm run production",
		)
	}

	if answers.Stack == models.Django {
		if managePyPath := utils.FindFile(
			path.Join(answers.WorkingDirectory, answers.ApplicationRoot),
			managePyFile,
		); managePyPath != "" {
			prefix := ""
			switch answers.DependencyManager {
			case models.Pipenv:
				prefix = "pipenv run "
			case models.Poetry:
				prefix = "poetry run "
			}

			managePyPath, _ = filepath.Rel(path.Join(answers.WorkingDirectory, answers.ApplicationRoot), managePyPath)
			answers.BuildSteps = append(
				answers.BuildSteps,
				"# Collect static files so that they can be served by Platform.sh",
				fmt.Sprintf("%spython %s collectstatic --noinput", prefix, managePyPath),
			)
		}
	}

	return nil
}
