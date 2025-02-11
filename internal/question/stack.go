package question

import (
	"context"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/colors"
	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/internal/questionnaire"
	"github.com/platformsh/platformify/internal/utils"
	"github.com/platformsh/platformify/platformifier"
	"github.com/platformsh/platformify/vendorization"
)

const (
	settingsPyFile   = "settings.py"
	managePyFile     = "manage.py"
	composerJSONFile = "composer.json"
	packageJSONFile  = "package.json"
	symfonyLockFile  = "symfony.lock"
	rackFile         = "config.ru"
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
	case platformifier.Symfony:
		// Pass to handle below if no interaction
		if answers.NoInteraction {
			answers.Stack = models.GenericStack
			return nil
		}
	default:
		answers.Stack = models.GenericStack
		return nil
	}

	rackPath := utils.FindFile(answers.WorkingDirectory, "", rackFile)
	if rackPath != "" {
		f, err := os.Open(rackPath)
		if err == nil {
			defer f.Close()
			if ok, _ = utils.ContainsStringInFile(f, "Rails.application.load_server", true); ok {
				answers.Stack = models.Rails
				return nil
			}
		}
	}

	requirementsPath := utils.FindFile(answers.WorkingDirectory, "", "requirements.txt")
	if requirementsPath != "" {
		f, err := answers.WorkingDirectory.Open(requirementsPath)
		if err == nil {
			defer f.Close()
			if ok, _ = utils.ContainsStringInFile(f, "flask", true); ok {
				answers.Stack = models.Flask
				return nil
			}
		}
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
	_, stderr, ok := colors.FromContext(ctx)
	if !ok {
		return questionnaire.ErrSilent
	}
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

	answers.Stack = models.GenericStack
	return nil
}
