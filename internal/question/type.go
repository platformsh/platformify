package question

import (
	"context"
	"fmt"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/colors"
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

	defer func() {
		_, stderr, ok := colors.FromContext(ctx)
		if !ok {
			return
		}

		if answers.Type.Runtime.Title() != "" {
			fmt.Fprintf(
				stderr,
				"%s %s\n",
				colors.Colorize(colors.GreenCode, "✓"),
				colors.Colorize(
					colors.BrandCode,
					fmt.Sprintf("Detected runtime: %s", answers.Type.Runtime.Title()),
				),
			)
		}
	}()

	typ, err := answers.Discoverer.Type()
	if err != nil {
		return err
	}
	runtime, _ := models.Runtimes.RuntimeByType(typ)

	if answers.NoInteraction {
		if runtime == nil {
			return fmt.Errorf("no runtime detected")
		}
		answers.Type.Runtime = *runtime
		answers.Type.Version = runtime.DefaultVersion()
		return nil
	}

	if runtime == nil || answers.Stack == models.GenericStack {
		question := &survey.Select{
			Message: "What language is your project using? We support the following:",
			Options: models.Runtimes.AllTitles(),
		}
		if runtime != nil {
			question.Default = runtime.Title()
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
	answers.Type.Runtime = *runtime
	answers.Type.Version = runtime.DefaultVersion()

	return nil
}
