package question

import (
	"context"

	"github.com/platformsh/platformify/internal/question/models"
)

type Environment struct{}

func (q *Environment) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	environment, err := answers.Discoverer.Environment()
	if err != nil {
		return err
	}

	answers.Environment = environment
	return nil
}
