package question

import (
	"context"
	"path/filepath"

	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/internal/utils"
)

type Locations struct{}

func (q *Locations) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}
	answers.Locations = answers.Type.Runtime.Docs.Locations
	if answers.Locations == nil {
		answers.Locations = make(map[string]map[string]interface{})
	}
	switch answers.Stack {
	case models.Django:
		answers.Locations["/static"] = map[string]interface{}{
			"root":    "static",
			"expires": "1h",
			"allow":   true,
		}
	case models.Symfony:
		answers.Locations["/"] = map[string]interface{}{
			"root":     "public",
			"passthru": "/index.php",
		}
	default:
		if answers.Type.Runtime.Type == "php" {
			locations := map[string]interface{}{
				"passthru": "/index.php",
				"root":     "",
			}
			if indexPath := utils.FindFile(answers.WorkingDirectory, "", "index.php"); indexPath != "" {
				if filepath.Dir(indexPath) != "." {
					locations["root"] = filepath.Dir(indexPath)
				}
			}
			answers.Locations["/"] = locations
		}
	}

	return nil
}
