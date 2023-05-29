package question

import (
	"context"

	"github.com/platformsh/platformify/internal/question/models"
)

type Locations struct{}

func (q *Locations) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}
	answers.Locations = make(map[string]map[string]interface{})
	if answers.Stack == models.Django {
		answers.Locations["/static"] = map[string]interface{}{
			"root":    "static",
			"expires": "1h",
			"allow":   true,
		}
	}

	return nil
}
