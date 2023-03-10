package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/answer"
)

type WebCommand struct{}

func (q *WebCommand) Ask(ctx context.Context) error {
	answers, ok := answer.FromContext(ctx)
	if !ok {
		return nil
	}

	question := &survey.Input{Message: "Web command:"}

	var webCommand string
	err := survey.AskOne(question, &webCommand)
	if err != nil {
		return err
	}

	answers.WebCommand = webCommand

	return nil
}
