package question

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/colors"
	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/internal/questionnaire"
	"github.com/platformsh/platformify/internal/utils"
)

const (
	settingsPyFile   = "settings.py"
	managePyFile     = "manage.py"
	composerJSONFile = "composer.json"
	packageJSONFile  = "package.json"
	symfonyLockFile  = "symfony.lock"
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

		if _, ok := utils.GetJSONKey([]string{"dependencies", "@strapi/strapi"}, packageJSONPath); ok {
			answers.Stack = models.Strapi
			return nil
		}
	}

	hasSymfonyLock := utils.FileExists(answers.WorkingDirectory, symfonyLockFile)
	hasSymfonyBundle := false
	hasIbexaDependencies := false
	hasShopwareDependencies := false
	for _, composerJSONPath := range composerJSONPaths {
		if _, ok := utils.GetJSONKey([]string{"autoload", "psr-0", "Shopware"}, composerJSONPath); ok {
			hasShopwareDependencies = true
			break
		}
		if _, ok := utils.GetJSONKey([]string{"autoload", "psr-4", "Shopware\\Core\\"}, composerJSONPath); ok {
			hasShopwareDependencies = true
			break
		}
		if _, ok := utils.GetJSONKey([]string{"autoload", "psr-4", "Shopware\\AppBundle\\"}, composerJSONPath); ok {
			hasShopwareDependencies = true
			break
		}

		if keywords, ok := utils.GetJSONKey([]string{"keywords"}, composerJSONPath); ok {
			if keywordsVal, ok := keywords.([]string); ok && slices.Contains(keywordsVal, "shopware") {
				hasShopwareDependencies = true
				break
			}
		}
		if requirements, ok := utils.GetJSONKey([]string{"require"}, composerJSONPath); ok {
			if requirementsVal, requirementsOK := requirements.(map[string]interface{}); requirementsOK {
				if _, hasSymfonyFrameworkBundle := requirementsVal["symfony/framework-bundle"]; hasSymfonyFrameworkBundle {
					hasSymfonyBundle = true
				}

				for requirement := range requirementsVal {
					if strings.HasPrefix(requirement, "shopware/") {
						hasShopwareDependencies = true
						break
					}
					if strings.HasPrefix(requirement, "ibexa/") {
						hasIbexaDependencies = true
						break
					}
					if strings.HasPrefix(requirement, "ezsystems/") {
						hasIbexaDependencies = true
						break
					}
				}
			}
		}
	}

	isSymfony := hasSymfonyBundle || hasSymfonyLock
	if isSymfony && !hasIbexaDependencies && !hasShopwareDependencies {
		_, stderr, ok := colors.FromContext(ctx)
		if !ok {
			return questionnaire.ErrSilent
		}

		confirm := true
		err := survey.AskOne(
			&survey.Confirm{
				Message: "It seems like this is a Symfony project, would you like to use the Symfony CLI to deploy your project instead?", //nolint:lll
				Default: confirm,
			},
			&confirm,
		)
		if err != nil {
			return err
		}

		if confirm {
			fmt.Fprintln(
				stderr,
				colors.Colorize(
					colors.WarningCode,
					"Check out the Symfony CLI documentation here: https://docs.platform.sh/guides/symfony/get-started.html",
				),
			)
			return questionnaire.ErrSilent
		}
	}

	return nil
}
