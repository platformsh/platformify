package question

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path"
	"strings"
)

type WorkingDirectory struct{}

func (q *WorkingDirectory) Ask(ctx context.Context) error {
	var outBuf, errBuf bytes.Buffer
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--git-dir")
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		if strings.Contains(errBuf.String(), "not a git repository") {
			return fmt.Errorf("current working directory is not a git repository")
		}
		return errors.New(errBuf.String())
	}

	gitRepoAbsPath := path.Dir(outBuf.String())
	if gitRepoAbsPath != "." {
		return fmt.Errorf("to run the command, go to the %s directory first", gitRepoAbsPath)
	}
	return nil
}
