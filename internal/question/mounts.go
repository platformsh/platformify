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
	case models.Rails:
		answers.Disk = "2048" // in MB
		answers.Mounts = map[string]map[string]string{
			"storage": {
				"source":      "local",
				"source_path": "storage",
			},
			"tmp": {
				"source":      "tmp",
				"source_path": "tmp",
			},
			"log": {
				"source":      "tmp",
				"source_path": "tmp",
			},
		}
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
	case models.Strapi:
		answers.Disk = "1024" // in MB
		answers.Mounts = map[string]map[string]string{
			"/.cache": {
				"source":      "local",
				"source_path": "cache",
			},
			"/.tmp": {
				"source":      "local",
				"source_path": "app",
			},
			"database": {
				"source":      "local",
				"source_path": "database",
			},
			"extensions": {
				"source":      "local",
				"source_path": "extensions",
			},
			"public/uploads": {
				"source":      "local",
				"source_path": "uploads",
			},
		}
	}

	return nil
}
