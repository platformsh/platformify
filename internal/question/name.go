package question

import (
	"context"
	"path"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/models"
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

	question := &survey.Input{Message: "Tell us your project name:", Default: path.Base(answers.WorkingDirectory)}

	var name string
	err := survey.AskOne(question, &name, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}

	answers.Name = name

	return nil
}
