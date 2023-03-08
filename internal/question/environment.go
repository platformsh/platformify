package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/answer"
)

type Environment struct{}

func (q *Environment) Ask(ctx context.Context) error {
	answers, ok := answer.FromContext(ctx)
	if !ok {
		return nil
	}
	defer func() {
		ctx = answer.ToContext(ctx, answers)
	}()

	var addVariable = true
	question := &survey.Confirm{
		Message: "Do you want to add a variable?",
		Default: true,
	}

	err := survey.AskOne(question, &addVariable)
	if err != nil {
		return err
	}

	if !addVariable {
		return nil
	}

	for addVariable {
		var qs = []*survey.Question{
			{
				Name:     "key",
				Prompt:   &survey.Input{Message: "Env var name"},
				Validate: survey.Required,
			},
			{
				Name:     "value",
				Prompt:   &survey.Input{Message: "Env var value"},
				Validate: survey.Required,
			},
		}

		variable := struct {
			Key   string
			Value string
		}{}

		err = survey.Ask(qs, &variable)
		if err != nil {
			return err
		}

		answers.Environment[variable.Key] = variable.Value

		question := &survey.Confirm{
			Message: "Do you want to add one more variable?",
			Default: true,
		}

		err := survey.AskOne(question, &addVariable)
		if err != nil {
			return err
		}
	}

	return nil
}
