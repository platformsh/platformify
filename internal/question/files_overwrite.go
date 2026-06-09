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

type FilesOverwrite struct {
	FilesToCreateUpdate []string
}

func (q *FilesOverwrite) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok || answers.NoInteraction {
		return nil
	}

	_, stderr, ok := colors.FromContext(ctx)
	if !ok {
		return nil
	}

	assets, _ := vendorization.FromContext(ctx)
	existingFiles := make([]string, 0, len(q.FilesToCreateUpdate))
	for _, p := range q.FilesToCreateUpdate {
		if st, err := fs.Stat(answers.WorkingDirectory, p); err == nil && !st.IsDir() {
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
