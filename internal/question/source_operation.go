package question

import (
	"context"

	"github.com/platformsh/platformify/internal/models"
)

type SourceOperation struct{}

func (q *SourceOperation) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}
	if len(answers.SourceOperations) != 0 {
		// Skip the step
		return nil
	}

	if answers.Stack == models.NextJS {
		answers.SourceOperations = map[string][]string{
			"auto-update": {
				"curl -fsS https://raw.githubusercontent.com/platformsh/source-operations/main/setup.sh | " +
					"{ bash /dev/fd/3 sop-autoupdate; } 3<&0",
			},
		}
	}

	return nil
}
