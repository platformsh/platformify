package question

import (
	"context"
	"fmt"
	"path"
	"path/filepath"

	"golang.org/x/exp/slices"

	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/internal/utils"
)

type BuildSteps struct{}

func (q *BuildSteps) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	for _, dm := range answers.DependencyManagers {
		switch dm {
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
			)
			if _, ok := utils.GetJSONValue(
				[]string{"scripts", "build"},
				path.Join(answers.WorkingDirectory, "package.json"),
				true,
			); ok {
				answers.BuildSteps = append(answers.BuildSteps, "yarn build")
			}
		case models.Npm:
			answers.BuildSteps = append(
				answers.BuildSteps,
				"npm i",
			)
			if _, ok := utils.GetJSONValue(
				[]string{"scripts", "build"},
				path.Join(answers.WorkingDirectory, "package.json"),
				true,
			); ok {
				answers.BuildSteps = append(answers.BuildSteps, "npm run build")
			}
		case models.Composer:
			answers.BuildSteps = append(
				answers.BuildSteps,
				"composer --no-ansi --no-interaction install --no-progress --prefer-dist --optimize-autoloader --no-dev",
			)
		}
	}

	if answers.Stack == models.Django {
		if managePyPath := utils.FindFile(
			path.Join(answers.WorkingDirectory, answers.ApplicationRoot),
			managePyFile,
		); managePyPath != "" {
			prefix := ""
			if slices.Contains(answers.DependencyManagers, models.Pipenv) {
				prefix = "pipenv run "
			} else if slices.Contains(answers.DependencyManagers, models.Poetry) {
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
