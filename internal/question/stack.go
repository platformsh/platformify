package question

import (
	"context"

	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/internal/utils"
)

const (
	settingsPyFile   = "settings.py"
	managePyFile     = "manage.py"
	composerJSONFile = "composer.json"
	packageJSONFile  = "package.json"
)

type Stack struct{}

func (q *Stack) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	answers.Stack = models.GenericStack

	hasSettingsPy := utils.FileExists(answers.WorkingDirectory, settingsPyFile)
	hasManagePy := utils.FileExists(answers.WorkingDirectory, managePyFile)
	if hasSettingsPy && hasManagePy {
		answers.Stack = models.Django
		return nil
	}

	composerJSONPaths := utils.FindAllFiles(answers.WorkingDirectory, composerJSONFile)
	for _, composerJSONPath := range composerJSONPaths {
		if _, ok := utils.GetJSONKey([]string{"require", "laravel/framework"}, composerJSONPath); ok {
			answers.Stack = models.Laravel
			return nil
		}
	}

	packageJSONPaths := utils.FindAllFiles(answers.WorkingDirectory, packageJSONFile)
	for _, packageJSONPath := range packageJSONPaths {
		if _, ok := utils.GetJSONKey([]string{"dependencies", "next"}, packageJSONPath); ok {
			answers.Stack = models.NextJS
			return nil
		}
	}

	return nil
}
