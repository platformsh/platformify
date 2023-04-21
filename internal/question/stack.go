package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/models"
	"github.com/platformsh/platformify/internal/utils"
)

type Stack struct{}

func (q *Stack) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	question := &survey.Select{
		Message: "What Stack is your project using?",
		Options: models.Stacks.AllTitles(),
	}
	hasSettingsPy := utils.FileExists(answers.WorkingDirectory, "settings.py")
	hasManagePy := utils.FileExists(answers.WorkingDirectory, "manage.py")
	if hasSettingsPy && hasManagePy {
		question.Default = models.Django.Title()
	}

	var title string
	err := survey.AskOne(question, &title, survey.WithPageSize(len(question.Options)))
	if err != nil {
		return err
	}
	stack, err := models.Stacks.StackByTitle(title)
	if err != nil {
		return err
	}

	answers.Stack = stack

	return nil
}
