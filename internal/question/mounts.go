package question

import (
	"context"

	"github.com/platformsh/platformify/internal/question/models"
)

type Mounts struct{}

func (q *Mounts) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	switch answers.Stack {
	case models.Laravel:
		answers.Disk = "2048" // in MB
		answers.Mounts = map[string]map[string]string{
			"storage": {
				"source":      "local",
				"source_path": "storage",
			},
			"bootstrap/cache": {
				"source":      "local",
				"source_path": "cache",
			},
			"/.config": {
				"source":      "local",
				"source_path": "config",
			},
		}
	case models.NextJS:
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
