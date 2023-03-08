package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"
)

type BuildSteps Question

func (q *BuildSteps) Ask(ctx context.Context) error {
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

		q.Answers.BuildSteps = append(q.Answers.BuildSteps, step)

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
