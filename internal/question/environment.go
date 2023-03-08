package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"
)

type Environment Question

func (q *Environment) Ask(ctx context.Context) error {
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

		q.Answers.Environment[variable.Key] = variable.Value

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
