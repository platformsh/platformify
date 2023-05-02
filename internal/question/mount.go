package question

import (
	"context"

	"github.com/platformsh/platformify/internal/question/models"
)

type Mount struct{}

func (q *Mount) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	if answers.Stack == models.NextJS {
		answers.Disk = "512" // in MB
		answers.Mounts = map[string]map[string]string{
			"/.npm": {
				"source":      "local",
				"source_path": "npm",
			},
		}
	}

	return nil
}
