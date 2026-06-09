package discovery

import (
	"slices"
	"strings"

	"github.com/platformsh/platformify/internal/utils"
	"github.com/platformsh/platformify/platformifier"
)

// Returns the stack, either from memory or by discovering it on the spot
func (d *Discoverer) Stack() (platformifier.Stack, error) {
	if stack, ok := d.memory["stack"]; ok {
		return stack.(platformifier.Stack), nil
	}

	stack, err := d.discoverStack()
	if err != nil {
		return platformifier.Generic, err
	}

	d.memory["stack"] = stack
	return stack, nil
}

func (d *Discoverer) discoverStack() (platformifier.Stack, error) {
	hasSettingsPy := utils.FileExists(d.fileSystem, "", settingsPyFile)
	hasManagePy := utils.FileExists(d.fileSystem, "", managePyFile)
	if hasSettingsPy && hasManagePy {
		return platformifier.Django, nil
	}

	requirementsPath := utils.FindFile(d.fileSystem, "", "requirements.txt")
	if requirementsPath != "" {
		f, err := d.fileSystem.Open(requirementsPath)
		if err == nil {
			defer f.Close()
			if ok, _ := utils.ContainsStringInFile(f, "flask", true); ok {
				return platformifier.Flask, nil
			}
		}
	}

	pyProjectPath := utils.FindFile(d.fileSystem, "", "pyproject.toml")
	if pyProjectPath != "" {
		if _, ok := utils.GetTOMLValue(
			d.fileSystem,
			[]string{"tool", "poetry", "dependencies", "flask"},
			pyProjectPath,
			true,
		); ok {
			return platformifier.Flask, nil
		}
	}

	pipfilePath := utils.FindFile(d.fileSystem, "", "Pipfile")
	if pipfilePath != "" {
		if _, ok := utils.GetTOMLValue(
			d.fileSystem,
			[]string{"packages", "flask"},
			pipfilePath,
			true,
		); ok {
			return platformifier.Flask, nil
		}
	}

	composerJSONPaths := utils.FindAllFiles(d.fileSystem, "", composerJSONFile)
	for _, composerJSONPath := range composerJSONPaths {
		if _, ok := utils.GetJSONValue(
			d.fileSystem,
			[]string{"require", "laravel/framework"},
			composerJSONPath,
			true,
		); ok {
			return platformifier.Laravel, nil
		}
	}

	packageJSONPaths := utils.FindAllFiles(d.fileSystem, "", packageJSONFile)
	for _, packageJSONPath := range packageJSONPaths {
		if _, ok := utils.GetJSONValue(
			d.fileSystem,
			[]string{"dependencies", "next"},
			packageJSONPath,
			true,
		); ok {
			return platformifier.NextJS, nil
		}

		if _, ok := utils.GetJSONValue(
			d.fileSystem,
			[]string{"dependencies", "@strapi/strapi"},
			packageJSONPath,
			true,
		); ok {
			return platformifier.Strapi, nil
		}

		if _, ok := utils.GetJSONValue(
			d.fileSystem,
			[]string{"dependencies", "strapi"},
			packageJSONPath,
			true,
		); ok {
			return platformifier.Strapi, nil
		}

		if _, ok := utils.GetJSONValue(
			d.fileSystem,
			[]string{"dependencies", "express"},
			packageJSONPath,
			true,
		); ok {
			return platformifier.Express, nil
		}
	}

	hasSymfonyLock := utils.FileExists(d.fileSystem, "", symfonyLockFile)
	hasSymfonyBundle := false
	for _, composerJSONPath := range composerJSONPaths {
		if _, ok := utils.GetJSONValue(
			d.fileSystem,
			[]string{"autoload", "psr-0", "shopware"},
			composerJSONPath,
			true,
		); ok {
			return platformifier.Shopware, nil
		}

		if _, ok := utils.GetJSONValue(
			d.fileSystem,
			[]string{"autoload", "psr-4", "shopware\\core\\"},
			composerJSONPath,
			true,
		); ok {
			return platformifier.Shopware, nil
		}

		if _, ok := utils.GetJSONValue(
			d.fileSystem,
			[]string{"autoload", "psr-4", "shopware\\appbundle\\"},
			composerJSONPath,
			true,
		); ok {
			return platformifier.Shopware, nil
		}

		if keywords, ok := utils.GetJSONValue(
			d.fileSystem,
			[]string{"keywords"},
			composerJSONPath,
			true,
		); ok {
			if keywordsVal, ok := keywords.([]string); ok && slices.Contains(keywordsVal, "shopware") {
				return platformifier.Shopware, nil
			}
		}
		if requirements, ok := utils.GetJSONValue(
			d.fileSystem,
			[]string{"require"},
			composerJSONPath,
			true,
		); ok {
			if requirementsVal, requirementsOK := requirements.(map[string]interface{}); requirementsOK {
				if _, requiresSymfony := requirementsVal["symfony/framework-bundle"]; requiresSymfony {
					hasSymfonyBundle = true
				}

				for requirement := range requirementsVal {
					if strings.HasPrefix(requirement, "shopware/") {
						return platformifier.Shopware, nil
					}
					if strings.HasPrefix(requirement, "ibexa/") {
						return platformifier.Ibexa, nil
					}
					if strings.HasPrefix(requirement, "ezsystems/") {
						return platformifier.Ibexa, nil
					}
				}
			}
		}
	}

	isSymfony := hasSymfonyBundle || hasSymfonyLock
	if isSymfony {
		return platformifier.Symfony, nil
	}

	return platformifier.Generic, nil
}
