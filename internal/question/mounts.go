package question

import (
	"context"

	"github.com/platformsh/platformify/internal/models"
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
			"storage/app/public": {
				"source":      "local",
				"source_path": "public",
			},
			"storage/framework/views": {
				"source":      "local",
				"source_path": "views",
			},
			"storage/framework/sessions": {
				"source":      "local",
				"source_path": "sessions",
			},
			"storage/framework/cache": {
				"source":      "local",
				"source_path": "cache",
			},
			"storage/logs": {
				"source":      "local",
				"source_path": "logs",
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
