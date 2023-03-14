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

	if answers.Type.Runtime.String() == "" {
		question := &survey.Select{
			Message: "Choose your PSH type:",
			Options: models.Runtimes.AllTitles(),
		}

		var title string
		err := survey.AskOne(question, &title, survey.WithPageSize(len(question.Options)))
		if err != nil {
			return err
		}

		runtime, err := models.Runtimes.RuntimeByTitle(title)
		if err != nil {
			return err
		}

		answers.Type.Runtime = runtime
	}

	runtime, err := models.Runtimes.RuntimeByName(answers.Type.Runtime)
	if err != nil {
		return err
	}

	versions, ok := models.Types[runtime]
	if !ok || len(versions) == 0 {
		return nil
	}

	question := &survey.Select{
		Message:  "Choose " + runtime.Title() + " version:",
		Options:  versions,
		PageSize: len(versions),
	}

	var version string
	err = survey.AskOne(question, &version)
	if err != nil {
		return err
	}

	answers.Type.Version = version

	return nil
}
