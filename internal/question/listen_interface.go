package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/models"
)

type ListenInterface struct{}

func (q *ListenInterface) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	interfaces := []string{
		"HTTP",
		"Unix-socket",
	}

	question := &survey.Select{
		Message: "Choose interface to listen to:",
		Options: interfaces,
	}

	var listen string
	err := survey.AskOne(question, &listen)
	if err != nil {
		return err
	}

	answers.ListenInterface = listen

	return nil
}
