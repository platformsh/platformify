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
	defer func() {
		ctx = answer.ToContext(ctx, answers)
	}()

	question := &survey.Input{Message: "What is the application root?"}

	var applicationRoot string
	err := survey.AskOne(question, &applicationRoot)
	if err != nil {
		return err
	}

	answers.ApplicationRoot = applicationRoot

	return nil
}
