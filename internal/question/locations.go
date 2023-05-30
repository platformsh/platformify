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
	answers.Locations = make(map[string]map[string]interface{})
	switch answers.Stack {
	case models.Django:
		answers.Locations["/static"] = map[string]interface{}{
			"root":    "static",
			"expires": "1h",
			"allow":   true,
		}
	default:
		if answers.Type.Runtime == models.PHP {
			locations := map[string]interface{}{
				"passthru": "/index.php",
			}
			if indexPath := utils.FindFile(answers.WorkingDirectory, "index.php"); indexPath != "" {
				indexRelPath, _ := filepath.Rel(answers.WorkingDirectory, indexPath)
				if filepath.Dir(indexRelPath) != "." {
					locations["root"] = filepath.Dir(indexRelPath)
				}
			}
			answers.Locations["/"] = locations
		}
	}

	return nil
}
