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
	defer func() {
		ctx = answer.ToContext(ctx, answers)
	}()

	var addStep = true
	question := &survey.Confirm{
		Message: "Do you want to add a build step?",
		Default: true,
	}

	err := survey.AskOne(question, &addStep)
	if err != nil {
		return err
	}

	if !addStep {
		return nil
	}

	for addStep {
		question2 := &survey.Input{Message: "Build step:"}

		var step string
		err = survey.AskOne(question2, &step)
		if err != nil {
			return err
		}

		answers.BuildSteps = append(answers.BuildSteps, step)

		question := &survey.Confirm{
			Message: "Do you want to add one more build step?",
			Default: true,
		}

		err := survey.AskOne(question, &addStep)
		if err != nil {
			return err
		}
	}

	return nil
}
