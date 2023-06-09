package question

import (
	"context"
	"fmt"
	"os"
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

	requirementsPath := utils.FindFile(answers.WorkingDirectory, "requirements.txt")
	if requirementsPath != "" {
		f, err := os.Open(requirementsPath)
		if err == nil {
			defer f.Close()
			if ok, _ := utils.ContainsStringInFile(f, "flask", true); ok {
				answers.Stack = models.Flask
				return nil
			}
		}
	}

	pyProjectPath := utils.FindFile(answers.WorkingDirectory, "pyproject.toml")
	if pyProjectPath != "" {
		if _, ok := utils.GetTOMLValue([]string{"tool", "poetry", "dependencies", "flask"}, pyProjectPath, true); ok {
			answers.Stack = models.Flask
			return nil
		}
	}

	pipfilePath := utils.FindFile(answers.WorkingDirectory, "Pipfile")
	if pipfilePath != "" {
		if _, ok := utils.GetTOMLValue([]string{"packages", "flask"}, pipfilePath, true); ok {
			answers.Stack = models.Flask
			return nil
		}
	}

	composerJSONPaths := utils.FindAllFiles(answers.WorkingDirectory, composerJSONFile)
	for _, composerJSONPath := range composerJSONPaths {
		if _, ok := utils.GetJSONValue([]string{"require", "laravel/framework"}, composerJSONPath, true); ok {
			answers.Stack = models.Laravel
			return nil
		}
	}

	packageJSONPaths := utils.FindAllFiles(answers.WorkingDirectory, packageJSONFile)
	for _, packageJSONPath := range packageJSONPaths {
		if _, ok := utils.GetJSONValue([]string{"dependencies", "next"}, packageJSONPath, true); ok {
			answers.Stack = models.NextJS
			return nil
		}

		if _, ok := utils.GetJSONValue([]string{"dependencies", "@strapi/strapi"}, packageJSONPath, true); ok {
			answers.Stack = models.Strapi
			return nil
		}
		if _, ok := utils.GetJSONValue([]string{"dependencies", "strapi"}, packageJSONPath, true); ok {
			answers.Stack = models.Strapi
			return nil
		}
	}

	hasSymfonyLock := utils.FileExists(answers.WorkingDirectory, symfonyLockFile)
	hasSymfonyBundle := false
	hasIbexaDependencies := false
	hasShopwareDependencies := false
	for _, composerJSONPath := range composerJSONPaths {
		if _, ok := utils.GetJSONValue([]string{"autoload", "psr-0", "shopware"}, composerJSONPath, true); ok {
			hasShopwareDependencies = true
			break
		}
		if _, ok := utils.GetJSONValue([]string{"autoload", "psr-4", "shopware\\core\\"}, composerJSONPath, true); ok {
			hasShopwareDependencies = true
			break
		}
		if _, ok := utils.GetJSONValue([]string{"autoload", "psr-4", "shopware\\appbundle\\"}, composerJSONPath, true); ok {
			hasShopwareDependencies = true
			break
		}

		if keywords, ok := utils.GetJSONValue([]string{"keywords"}, composerJSONPath, true); ok {
			if keywordsVal, ok := keywords.([]string); ok && slices.Contains(keywordsVal, "shopware") {
				hasShopwareDependencies = true
				break
			}
		}
		if requirements, ok := utils.GetJSONValue([]string{"require"}, composerJSONPath, true); ok {
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
				Message: "It seems like this project uses Symfony full-stack. For a better experience, you should use Symfony CLI. Would you like to use it to deploy your project instead?", //nolint:lll
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
