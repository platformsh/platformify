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

	rt := answers.Type.Runtime.Type
	switch rt {
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
