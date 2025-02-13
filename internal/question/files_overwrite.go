package question

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/colors"
	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/internal/questionnaire"
	"github.com/platformsh/platformify/vendorization"
)

// FilesOverwrite prompts the user to confirm overwriting existing config files.
// If FilesToCreateUpdate is set, those files are checked instead of the
// default proprietary files list.
type FilesOverwrite struct {
	FilesToCreateUpdate []string
}

func (q *FilesOverwrite) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	if answers.NoInteraction {
		return nil
	}

	_, stderr, ok := colors.FromContext(ctx)
	if !ok {
		return nil
	}

	filesToCheck := q.FilesToCreateUpdate
	if len(filesToCheck) == 0 {
		assets, _ := vendorization.FromContext(ctx)
		filesToCheck = assets.ProprietaryFiles()
	}

	existingFiles := make([]string, 0, len(filesToCheck))
	for _, p := range filesToCheck {
		if st, err := fs.Stat(answers.WorkingDirectory, p); err == nil && !st.IsDir() {
			existingFiles = append(existingFiles, p)
		}
	}

	if len(existingFiles) > 0 {
		assets, _ := vendorization.FromContext(ctx)
		fmt.Fprintln(
			stderr,
			colors.Colorize(
				colors.WarningCode,
				fmt.Sprintf("You are reconfiguring the project at %s.", answers.Cwd),
			),
		)
		fmt.Fprintln(
			stderr,
			colors.Colorize(
				colors.WarningCode,
				fmt.Sprintf(
					"The following %s files already exist in this directory:",
					assets.ServiceName,
				),
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
