package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/question/models"
)

type Type struct{}

func (q *Type) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}
	if answers.Type.String() != "" {
		// Skip the step
		return nil
	}

	runtime := models.RuntimeForStack(answers.Stack)
	if runtime == "" {
		question := &survey.Select{
			Message: "What language is your project using? We support the following:",
			Options: models.Runtimes.AllTitles(),
		}

		var title string
		err := survey.AskOne(question, &title, survey.WithPageSize(len(question.Options)))
		if err != nil {
			return err
		}

		runtime, err = models.Runtimes.RuntimeByTitle(title)
		if err != nil {
			return err
		}
	}
	answers.Type.Runtime = runtime

	// @todo The user may want a specific version. Should we ask instead of assuming?
	answers.Type.Version = models.DefaultVersionForRuntime(runtime)

	return nil
}
