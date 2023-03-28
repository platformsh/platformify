package question

import (
	"context"

	"github.com/platformsh/platformify/internal/models"
)

type Mount struct{}

func (q *Mount) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}
	if len(answers.Mounts) != 0 {
		// Skip the step
		return nil
	}

	if answers.Stack == models.NextJS {
		answers.Disk = "512" // in MB
		answers.Mounts = map[string]map[string]string{
			"/.next": {
				"source":      "local",
				"source_path": "next",
			},
		}
	}

	return nil
}
