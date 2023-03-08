package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/answer"
)

type Listen struct{}

func (q *Listen) Ask(ctx context.Context) error {
	answers, ok := answer.FromContext(ctx)
	if !ok {
		return nil
	}
	defer func() {
		ctx = answer.ToContext(ctx, answers)
	}()

	interfaces := []string{
		"HTTP",
		"Unix-socket",
	}

	question := &survey.Select{
		Message: "Choose interface to listen to:",
		Options: interfaces,
		Default: nil,
	}

	var listen string
	err := survey.AskOne(question, &listen)
	if err != nil {
		return err
	}

	answers.Listen = listen

	return nil
}
