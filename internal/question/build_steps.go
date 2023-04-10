package question

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/platformsh/platformify/internal/models"
	"github.com/platformsh/platformify/internal/utils"
)

const (
	managePyFile = "manage.py"
)

type BuildSteps struct{}

func (q *BuildSteps) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	switch answers.Stack {
	case models.Django:
		prefix := ""
		switch answers.DependencyManager {
		case models.Poetry:
			answers.BuildSteps = append(
				answers.BuildSteps,
				"# Set PIP_USER to 0 so that poetry does not complain",
				"export PIP_USER=0",
				"# Install poetry as a global tool",
				"python -m venv /app/.global",
				"pip install poetry==$POETRY_VERSION",
				"poetry install",
			)
			prefix = "poetry run "
		case models.Pipenv:
			answers.BuildSteps = append(
				answers.BuildSteps,
				"# Set PIP_USER to 0 so that Pipenv does not complain",
				"export PIP_USER=0",
				"# Install Pipenv as a global tool",
				"python -m venv /app/.global",
				"pip install poetry==$PIPENV_VERSION",
				"pipenv install",
			)
			prefix = "pipenv run "
		case models.Pip:
			answers.BuildSteps = append(
				answers.BuildSteps,
				"pip install -r requirements.txt",
			)
		}

		cwd, _ := os.Getwd()
		if managePyPath := utils.FindFile(path.Join(cwd, answers.ApplicationRoot), managePyFile); managePyPath != "" {
			managePyPath, _ = filepath.Rel(path.Join(cwd, answers.ApplicationRoot), managePyPath)
			answers.BuildSteps = append(
				answers.BuildSteps,
				"# Collect static files so that they can be served by Platform.sh",
				fmt.Sprintf("%spython %s collectstatic --noinput", prefix, managePyPath),
			)
		}
	case models.NextJS:
		switch answers.DependencyManager {
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
		}
	case models.Laravel:
		answers.BuildSteps = append(
			answers.BuildSteps,
			"composer --no-ansi --no-interaction install --no-progress --prefer-dist --optimize-autoloader --no-dev",
			"# Install a specific NodeJS version https://github.com/platformsh/snippets/",
			"# uncomment next line to build assets deploying",
			"# npm install && npm run production",
		)
	}

	if len(answers.BuildSteps) > 0 {
		fmt.Println("We identified a few build steps for you already!")
		for _, step := range answers.BuildSteps {
			fmt.Println("  " + step)
		}
	}

	return nil
}
