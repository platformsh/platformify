package question

import (
	"context"

	"github.com/platformsh/platformify/internal/models"
)

type SocketFamily struct{}

func (q *SocketFamily) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	switch answers.Type.Runtime {
	case models.PHP:
		return nil
	case models.Ruby, models.Python:
		answers.SocketFamily = models.UnixSocket
		return nil
	default:
		answers.SocketFamily = models.TCP
		return nil
	}
}
