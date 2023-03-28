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
	if answers.ListenInterface != "" {
		// Skip the step
		return nil
	}

	if answers.Stack == models.NextJS {
		return nil
	}

	question := &survey.Select{
		Message: "Choose interface to listen to:",
		Options: models.ListenInterfaces.AllTitles(),
	}

	var title string
	err := survey.AskOne(question, &title)
	if err != nil {
		return err
	}
	iface, err := models.ListenInterfaces.ListenInterfaceByTitle(title)
	if err != nil {
		return err
	}

	answers.ListenInterface = iface

	return nil
}
