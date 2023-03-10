package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/answer"
)

type BuildSteps struct{}

func (q *BuildSteps) Ask(ctx context.Context) error {
	answers, ok := answer.FromContext(ctx)
	if !ok {
		return nil
	}

	for {
		var question survey.Prompt
		var addStep = true
		question = &survey.Confirm{
			Message: "Do you want to add a build step?",
			Default: true,
		}

		err := survey.AskOne(question, &addStep)
		if err != nil {
			return err
		}

		if !addStep {
			break
		}

		question = &survey.Input{Message: "Build step:"}

		var step string
		err = survey.AskOne(question, &step)
		if err != nil {
			return err
		}

		answers.BuildSteps = append(answers.BuildSteps, step)
	}

	return nil
}
