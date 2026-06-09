package discovery

import (
	"fmt"
	"slices"

	"github.com/platformsh/platformify/internal/utils"
	"github.com/platformsh/platformify/platformifier"
)

// Returns the application build steps, either from memory or by discovering it on the spot
func (d *Discoverer) BuildSteps() ([]string, error) {
	if buildSteps, ok := d.memory["build_steps"]; ok {
		return buildSteps.([]string), nil
	}

	buildSteps, err := d.discoverBuildSteps()
	if err != nil {
		return nil, err
	}

	d.memory["build_steps"] = buildSteps
	return buildSteps, nil
}

func (d *Discoverer) discoverBuildSteps() ([]string, error) {
	dependencyManagers, err := d.DependencyManagers()
	if err != nil {
		return nil, err
	}

	typ, err := d.Type()
	if err != nil {
		return nil, err
	}

	stack, err := d.Stack()
	if err != nil {
		return nil, err
	}

	appRoot, err := d.ApplicationRoot()
	if err != nil {
		return nil, err
	}

	buildSteps := make([]string, 0)

	// Start with lower priority dependency managers first
	slices.Reverse(dependencyManagers)
	for _, dm := range dependencyManagers {
		switch dm {
		case "poetry":
			buildSteps = append(
				buildSteps,
				"# Set PIP_USER to 0 so that Poetry does not complain",
				"export PIP_USER=0",
				"# Install poetry as a global tool",
				"python -m venv /app/.global",
				"pip install poetry==$POETRY_VERSION",
				"poetry install",
			)
		case "pipenv":
			buildSteps = append(
				buildSteps,
				"# Set PIP_USER to 0 so that Pipenv does not complain",
				"export PIP_USER=0",
				"# Install Pipenv as a global tool",
				"python -m venv /app/.global",
				"pip install pipenv==$PIPENV_TOOL_VERSION",
				"pipenv install",
			)
		case "pip":
			buildSteps = append(
				buildSteps,
				"pip install -r requirements.txt",
			)
		case "yarn", "npm":
			// Install n, if on different runtime
			if typ != "nodejs" {
				buildSteps = append(
					buildSteps,
					"n auto || n lts",
					"hash -r",
				)
			}

			if dm == "yarn" {
				buildSteps = append(
					buildSteps,
					"yarn",
				)
			} else {
				buildSteps = append(
					buildSteps,
					"npm i",
				)
			}

			if _, ok := utils.GetJSONValue(
				d.fileSystem,
				[]string{"scripts", "build"},
				"package.json",
				true,
			); ok {
				buildSteps = append(buildSteps, d.nodeScriptPrefix()+"build")
			}
		case "composer":
			buildSteps = append(
				buildSteps,
				"composer --no-ansi --no-interaction install --no-progress --prefer-dist --optimize-autoloader --no-dev",
			)
		}
	}

	switch stack {
	case platformifier.Django:
		if managePyPath := utils.FindFile(
			d.fileSystem,
			appRoot,
			managePyFile,
		); managePyPath != "" {
			buildSteps = append(
				buildSteps,
				"# Collect static files",
				fmt.Sprintf("%spython %s collectstatic --noinput", d.pythonPrefix(), managePyPath),
			)
		}
	case platformifier.NextJS:
		// If there is no custom build script, fallback to next build for Next.js projects
		if !slices.Contains(buildSteps, "yarn build") && !slices.Contains(buildSteps, "npm run build") {
			buildSteps = append(buildSteps, d.nodeExecPrefix()+"next build")
		}
	}

	return buildSteps, nil
}
