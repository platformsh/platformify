package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/answer"
)

type ApplicationRoot struct{}

func (q *ApplicationRoot) Ask(ctx context.Context) error {
	answers, ok := answer.FromContext(ctx)
	if !ok {
		return nil
	}

	question := &survey.Input{Message: "Application root:"}

	var applicationRoot string
	err := survey.AskOne(question, &applicationRoot, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}

	answers.ApplicationRoot = applicationRoot

	return nil
}
