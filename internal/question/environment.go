package question

import (
	"context"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/models"
)

type Environment struct{}

func (q *Environment) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	for {
		question := &survey.Confirm{
			Message: "Do you want to add a variable?",
			Default: true,
		}

		var addVariable = true
		err := survey.AskOne(question, &addVariable)
		if err != nil {
			return err
		}

		if !addVariable {
			break
		}

		questions := []*survey.Question{
			{
				Name:      "key",
				Prompt:    &survey.Input{Message: "Env var name"},
				Validate:  survey.Required,
				Transform: survey.TransformString(strings.ToUpper),
			},
			{
				Name:   "value",
				Prompt: &survey.Input{Message: "Env var value"},
			},
		}

		variable := struct {
			Key   string
			Value string
		}{}

		err = survey.Ask(questions, &variable)
		if err != nil {
			return err
		}

		answers.Environment[variable.Key] = variable.Value
	}

	return nil
}
