package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/answer"
)

type Root struct{}

func (q *Root) Ask(ctx context.Context) error {
	answers, ok := answer.FromContext(ctx)
	if !ok {
		return nil
	}
	defer func() {
		ctx = answer.ToContext(ctx, answers)
	}()

	question := &survey.Input{Message: "What is the application root?"}

	var root string
	err := survey.AskOne(question, &root)
	if err != nil {
		return err
	}

	answers.Root = root

	return nil
}
