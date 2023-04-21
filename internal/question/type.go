package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/models"
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

	var runtime models.Runtime
	switch answers.Stack {
	case models.Django:
		runtime = models.Python
	case models.Laravel:
		runtime = models.PHP
	case models.NextJS:
		runtime = models.NodeJS
	default:
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

	versions, ok := models.LanguageTypeVersions[runtime]
	if !ok || len(versions) == 0 {
		return nil
	}
	answers.Type.Version = versions[0]

	return nil
}
