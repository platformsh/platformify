package question

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/models"
)

const (
	configReaderPkgName        = "platformsh-config"
	configReaderPkgDescription = "Platform.sh Config Reader"
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
			Message: fmt.Sprintf("Would you like to install %s (%s)?", configReaderPkgName, configReaderPkgDescription),
			Default: true,
		}

		var installPackage bool
		err := survey.AskOne(question, &installPackage)
		if err != nil {
			return err
		}

		if installPackage {
			var command string
			var args []string
			switch answers.DependencyManager {
			case models.Yarn:
				command = "yarn"
				args = []string{"add", configReaderPkgName}
			case models.Npm:
				command = "npm"
				args = []string{"install", configReaderPkgName, "--save"}
			}

			cmd := exec.CommandContext(ctx, command, args...)
			err := cmd.Run()
			if err != nil {
				return err
			}
			fmt.Printf("%s has been successfully installed.\n", configReaderPkgName)
		}
	default:
		// Skip the step
		return nil
	}

	return nil
}
