package questionnaire

import (
	"errors"
	"fmt"
)

var (
	ErrSilent      = errors.New("silent error")
	ErrUserAborted = fmt.Errorf("user aborted: %w", ErrSilent)
)
