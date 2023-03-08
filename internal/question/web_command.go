package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"
)

type WebCommand Question

func (q *WebCommand) Ask(ctx context.Context) error {
	question := &survey.Input{Message: "Web command?"}

	var webCommand string
	err := survey.AskOne(question, &webCommand)
	if err != nil {
		return err
	}

	q.Answers.WebCommand = webCommand

	return nil
}
