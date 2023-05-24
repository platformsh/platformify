package question

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/colors"
	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/internal/questionnaire"
)

var proprietaryFiles = []string{
	".environment",
	".platform.app.yaml",
	".platform/services.yaml",
	".platform/routes.yaml",
	".platform/applications.yaml",
}

type FilesOverwrite struct{}

func (q *FilesOverwrite) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	_, stderr, ok := colors.FromContext(ctx)
	if !ok {
		return nil
	}

	existingFiles := make([]string, 0, len(proprietaryFiles))
	for _, p := range proprietaryFiles {
		if st, err := os.Stat(filepath.Join(answers.WorkingDirectory, p)); err == nil && !st.IsDir() {
			existingFiles = append(existingFiles, p)
		}
	}

	if len(existingFiles) > 0 {
		fmt.Fprintln(
			stderr,
			colors.Colorize(
				colors.WarningCode,
				fmt.Sprintf("You are reconfiguring the project at %s.", answers.WorkingDirectory),
			),
		)
		fmt.Fprintln(
			stderr,
			colors.Colorize(
				colors.WarningCode,
				"The following Platform.sh files already exist in this directory:",
			),
		)
		for _, p := range existingFiles {
			fmt.Fprintln(stderr, colors.Colorize(colors.WarningCode, fmt.Sprintf("  - %s", p)))
		}
		proceed := false
		if err := survey.AskOne(&survey.Confirm{
			Message: "Do you want to overwrite them?",
			Default: proceed,
		}, &proceed); err != nil {
			return err
		}

		if !proceed {
			return questionnaire.ErrUserAborted
		}
	}

	return nil
}
