package question

import (
	"context"
	"fmt"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/colors"
	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/internal/questionnaire"
	"github.com/platformsh/platformify/platformifier"
	"github.com/platformsh/platformify/vendorization"
)

type Stack struct{}

func (q *Stack) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	defer func() {
		_, stderr, stderrOK := colors.FromContext(ctx)
		if !stderrOK {
			return
		}
		if answers.Stack != models.GenericStack {
			fmt.Fprintf(
				stderr,
				"%s %s\n",
				colors.Colorize(colors.GreenCode, "✓"),
				colors.Colorize(
					colors.BrandCode,
					fmt.Sprintf("Detected stack: %s", answers.Stack.Title()),
				),
			)
		}
	}()

	answers.Stack = models.GenericStack
	stack, err := answers.Discoverer.Stack()
	if err != nil {
		return err
	}

	switch stack {
	case platformifier.Django:
		answers.Stack = models.Django
		return nil
	case platformifier.Express:
		answers.Stack = models.Express
		return nil
	case platformifier.Flask:
		answers.Stack = models.Flask
		return nil
	case platformifier.Laravel:
		answers.Stack = models.Laravel
		return nil
	case platformifier.NextJS:
		answers.Stack = models.NextJS
		return nil
	case platformifier.Strapi:
		answers.Stack = models.Strapi
		return nil
	case platformifier.Rails:
		answers.Stack = models.Rails
		return nil
	case platformifier.Symfony:
		// Interactive: offer Symfony CLI below.
	default:
		answers.Stack = models.GenericStack
		return nil
	}

	_, stderr, ok := colors.FromContext(ctx)
	if !ok {
		return questionnaire.ErrSilent
	}

	confirm := true
	if err := survey.AskOne(
		&survey.Confirm{
			Message: "It seems like this project uses Symfony full-stack. For a better experience, you should use Symfony CLI. Would you like to use it to deploy your project instead?", //nolint:lll
			Default: confirm,
		},
		&confirm,
	); err != nil {
		return err
	}

	assets, _ := vendorization.FromContext(ctx)
	if confirm {
		fmt.Fprintln(
			stderr,
			colors.Colorize(
				colors.WarningCode,
				fmt.Sprintf(
					"Check out the Symfony CLI documentation here: %s",
					assets.Docs().SymfonyCLI,
				),
			),
		)
		return questionnaire.ErrSilent
	}

	return nil
}
