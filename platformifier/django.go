package platformifier

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/platformsh/platformify/internal/colors"
	"github.com/platformsh/platformify/internal/utils"
)

const (
	settingsPyFile        = "settings.py"
	settingsPshPyFile     = "settings_psh.py"
	importSettingsPshLine = "from settings_psh import *"
	djangoTemplatesPath   = "templates/django"
)

func newDjangoPlatformifier(templates fs.FS, file fileCreator) *djangoPlatformifier {
	return &djangoPlatformifier{
		templates: templates,
		file:      file,
	}
}

type djangoPlatformifier struct {
	templates fs.FS
	file      fileCreator
}

func (p *djangoPlatformifier) Platformify(ctx context.Context, input *UserInput) error {
	// Get working directory.
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current working directory: %w", err)
	}

	appRoot := path.Join(cwd, input.Root, input.ApplicationRoot)
	if settingsPath := utils.FindFile(appRoot, settingsPyFile); settingsPath != "" {
		pshSettingsPath := filepath.Join(filepath.Dir(settingsPath), settingsPshPyFile)
		tpl, parseErr := template.New(settingsPshPyFile).Funcs(sprig.FuncMap()).
			ParseFS(p.templates, settingsPshPyFile)
		if parseErr != nil {
			return fmt.Errorf("could not parse template: %w", parseErr)
		}
		f, err := p.file.Create(pshSettingsPath)
		if err != nil {
			return err
		}
		defer f.Close()

		err = tpl.Execute(f, input)
		if err != nil {
			return err
		}
	}

	// append from settings_psh import * to the bottom of settings.py
	if settingsPath := utils.FindFile(appRoot, settingsPyFile); settingsPath != "" {
		f, err := os.OpenFile(settingsPath, os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return nil
		}
		defer f.Close()

		// Check if there is an import line in the file
		found, err := utils.ContainsStringInFile(settingsPath, importSettingsPshLine)
		if err != nil {
			return err
		}

		if !found {
			if _, err = f.WriteString("\n\n" + importSettingsPshLine + "\n"); err != nil {
				out, _, ok := colors.FromContext(ctx)
				if !ok {
					return nil
				}

				fmt.Fprintf(
					out,
					colors.Colorize(
						colors.WarningCode,
						"We have created a %s file for you. Please add the following line to your %s file:\n",
					),
					settingsPshPyFile,
					settingsPyFile,
				)
				fmt.Fprint(out, colors.Colorize(colors.WarningCode, "    "+importSettingsPshLine+"\n"))
				return nil
			}
		}
	}

	return nil
}

// FIXME: !!!!!!!!!!!!!!!!!!!!!!
// func (p *djangoPlatformifier) Platformify2(ctx context.Context, input *UserInput) error {
// 	// Gather templates.
// 	templates, err := utils.GatherTemplates(ctx, p.templates, ".")
// 	if err != nil {
// 		return err
// 	}
//
// 	appRoot := path.Join(cwd, p.Root, p.ApplicationRoot)
// 	if settingsPath := utils.FindFile(appRoot, settingsPyFile); settingsPath != "" {
// 		tpl, parseErr := template.New(settingsPshPyFile).Funcs(sprig.FuncMap()).ParseFS(
// 			templatesFs, "templates/_extras/django/settings_psh.py",
// 		)
// 		if parseErr != nil {
// 			return fmt.Errorf("could not parse template: %w", parseErr)
// 		}
// 		pshSettingsPath, _ := filepath.Rel(
// 			cwd,
// 			filepath.Join(filepath.Dir(settingsPath), settingsPshPyFile),
// 		)
// 		templates[pshSettingsPath] = tpl
// 	}
//
// 	if err := utils.WriteTemplates(ctx, cwd, templates, p.UserInput); err != nil {
// 		return fmt.Errorf("could not write Platform.sh files: %w", err)
// 	}
//
// 	// append from settings_psh import * to the bottom of settings.py
// 	if settingsPath := utils.FindFile(appRoot, settingsPyFile); settingsPath != "" {
// 		f, err := os.OpenFile(settingsPath, os.O_APPEND|os.O_WRONLY, 0o644)
// 		if err != nil {
// 			return nil
// 		}
// 		defer f.Close()
//
// 		// Check if there is an import line in the file
// 		found, err := utils.ContainsStringInFile(settingsPath, importSettingsPshLine)
// 		if err != nil {
// 			return err
// 		}
//
// 		if !found {
// 			if _, err = f.WriteString("\n\n" + importSettingsPshLine + "\n"); err != nil {
// 				out, _, ok := colors.FromContext(ctx)
// 				if !ok {
// 					return nil
// 				}
//
// 				fmt.Fprintf(
// 					out,
// 					colors.Colorize(
// 						colors.WarningCode,
// 						"We have created a %s file for you. Please add the following line to your %s file:\n",
// 					),
// 					settingsPshPyFile,
// 					settingsPyFile,
// 				)
// 				fmt.Fprint(out, colors.Colorize(colors.WarningCode, "    "+importSettingsPshLine+"\n"))
// 				return nil
// 			}
// 		}
// 	}
//
// 	return nil
// }
