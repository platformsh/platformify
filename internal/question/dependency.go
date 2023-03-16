package question

import (
	"context"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/models"
)

type Dependency struct{}

func (q *Dependency) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	switch answers.DependencyManager {
	case models.Yarn, models.Npm:
		question := &survey.Confirm{
			Message: "Would you like to install platformsh-config (Platform.sh Config Reader)?",
			Default: true,
		}

		var installPackage bool
		err := survey.AskOne(question, &installPackage)
		if err != nil {
			return err
		}

		if installPackage {
			var command string
			switch answers.DependencyManager {
			case models.Yarn:
				command = "yarn add platformsh-config"
			case models.Npm:
				command = "npm install platformsh-config --save"
			}

			cmd := exec.CommandContext(ctx, command)
			err := cmd.Run()
			if err != nil {
				return nil
				// return err
			}
		}
	default:
		// Skip the step
		return nil
	}

	return nil
}
