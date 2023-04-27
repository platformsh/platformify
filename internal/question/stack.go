package question

import (
	"context"

	"github.com/platformsh/platformify/internal/models"
	"github.com/platformsh/platformify/internal/utils"
)

type Stack struct{}

func (q *Stack) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	answers.Stack = models.GenericStack

	hasSettingsPy := utils.FileExists(answers.WorkingDirectory, "settings.py")
	hasManagePy := utils.FileExists(answers.WorkingDirectory, "manage.py")
	if hasSettingsPy && hasManagePy {
		answers.Stack = models.Django
		return nil
	}

	composerJSONPaths := utils.FindAllFiles(answers.WorkingDirectory, "composer.json")
	for _, composerJSONPath := range composerJSONPaths {
		if _, ok := utils.GetJSONKey([]string{"require", "laravel/framework"}, composerJSONPath); ok {
			answers.Stack = models.Laravel
			return nil
		}
	}

	packageJSONPaths := utils.FindAllFiles(answers.WorkingDirectory, "package.json")
	for _, packageJSONPath := range packageJSONPaths {
		if _, ok := utils.GetJSONKey([]string{"dependencies", "next"}, packageJSONPath); ok {
			answers.Stack = models.NextJS
			return nil
		}
	}

	return nil
}
