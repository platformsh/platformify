package question

import (
	"context"

	"github.com/platformsh/platformify/internal/models"
)

type ListenInterface struct{}

func (q *ListenInterface) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	switch answers.Type.Runtime.String() {
	case models.PHP.String():
		return nil
	case models.Ruby.String(), models.Python.String():
		answers.ListenInterface = models.UnixSocket
		return nil
	default:
		answers.ListenInterface = models.HTTP
		return nil
	}
}
