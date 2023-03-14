package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/models"
)

type Stack struct{}

func (q *Stack) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	question := &survey.Select{
		Message: "Choose your stack:",
		Options: models.Stacks.AllTitles(),
	}

	var title string
	err := survey.AskOne(question, &title, survey.WithPageSize(len(question.Options)))
	if err != nil {
		return err
	}
	stack, err := models.Stacks.StackByTitle(title)
	if err != nil {
		return err
	}

	answers.Stack = stack
	if stack != models.GenericStack {
		var runtime models.Runtime
		switch stack {
		case models.Django:
			runtime = models.Python
		case models.Laravel:
			runtime = models.PHP
		case models.NextJS:
			runtime = models.NodeJS
		}

		answers.Type.Runtime = runtime
	}

	return nil
}
