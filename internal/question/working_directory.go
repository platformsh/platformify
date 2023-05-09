package question

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/platformsh/platformify/internal/colors"
	"github.com/platformsh/platformify/internal/models"
)

type WorkingDirectory struct{}

func (q *WorkingDirectory) Ask(ctx context.Context) error {
	_, stderr, ok := colors.FromContext(ctx)
	if !ok {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}
	answers.WorkingDirectory = cwd
	answers.HasGit = false

	var outBuf, errBuf bytes.Buffer
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--git-dir")
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err = cmd.Run()
	if err != nil {
		fmt.Fprintln(
			stderr,
			colors.Colorize(
				colors.WarningCode,
				"It seems like you're not inside a Git repository.",
			),
		)
		fmt.Fprintln(
			stderr,
			colors.Colorize(
				colors.WarningCode,
				"You'll need to create a Git repository to deploy your application to Platform.sh",
			),
		)
		return nil
	}

	gitRepoAbsPath := path.Dir(outBuf.String())
	if gitRepoAbsPath != "." {
		return fmt.Errorf("platformify should be run at the root of a Git repository, "+
			"please change the directory to %s and run the command again", gitRepoAbsPath)
	}

	return nil
}
