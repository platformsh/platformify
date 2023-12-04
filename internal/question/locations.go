package question

import (
	"context"
	"path"

	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/internal/utils"
)

type Locations struct{}

func (q *Locations) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}
	answers.Locations = make(map[string]map[string]interface{})
	switch answers.Stack {
	case models.Django:
		answers.Locations["/static"] = map[string]interface{}{
			"root":    "static",
			"expires": "1h",
			"allow":   true,
		}
	default:
		if answers.Type.Runtime.Type == "php" {
			locations := map[string]interface{}{
				"passthru": "/index.php",
				"root":     "",
			}
			if indexPath := utils.FindFile(answers.WorkingDirectory, "", "index.php"); indexPath != "" {
				if path.Dir(indexPath) != "." {
					locations["root"] = path.Dir(indexPath)
				}
			}
			answers.Locations["/"] = locations
		}
	}

	return nil
}
