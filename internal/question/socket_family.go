package question

import (
	"context"

	"github.com/platformsh/platformify/internal/question/models"
)

type SocketFamily struct{}

func (q *SocketFamily) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	switch answers.Type.Runtime.Type {
	case "php":
		return nil
	case "ruby", "python":
		answers.SocketFamily = models.UnixSocket
		return nil
	default:
		answers.SocketFamily = models.TCP
		return nil
	}
}
