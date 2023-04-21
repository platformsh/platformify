package question

import (
	"context"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/internal/utils"
)

const (
	settingsPyFile = "settings.py"
	managePyFile   = "manage.py"
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
	hasSettingsPy := utils.FileExists(answers.WorkingDirectory, settingsPyFile)
	hasManagePy := utils.FileExists(answers.WorkingDirectory, managePyFile)
	if hasSettingsPy && hasManagePy {
		question.Default = models.Django.Title()
	}

	var stack models.Stack
	err := survey.AskOne(question, &stack, survey.WithPageSize(len(question.Options)))
	if err != nil {
		return err
	}

	answers.Stack = stack

	return nil
}
