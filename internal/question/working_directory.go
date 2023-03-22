package question

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"path"
)

type WorkingDirectory struct{}

func (q *WorkingDirectory) Ask(ctx context.Context) error {
	var outBuf, errBuf bytes.Buffer
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--git-dir")
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		// TODO: print the following log message only when we are in debug mode.
		log.Println(errBuf.String())
		return fmt.Errorf("platformify should be run at the root of a Git repository, " +
			"please change to a Git directory and run the command again")
	}

	gitRepoAbsPath := path.Dir(outBuf.String())
	if gitRepoAbsPath != "." {
		return fmt.Errorf("platformify should be run at the root of a Git repository, "+
			"please change the directory to %s and run the command again", gitRepoAbsPath)
	}
	return nil
}
