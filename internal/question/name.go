package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/question/models"
)

type Name struct{}

func (q *Name) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}
	if answers.Name != "" {
		// Skip the step
		return nil
	}

	question := &survey.Input{Message: "Application name:", Default: "app"}

	var name string
	err := survey.AskOne(question, &name, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}

	answers.Name = name

	return nil
}
