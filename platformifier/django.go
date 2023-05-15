package platformifier

import (
	"context"
	"fmt"
	"io/fs"
	"os"
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
)

func newDjangoPlatformifier(templates fs.FS, fileSystem FS) *djangoPlatformifier {
	return &djangoPlatformifier{
		templates:  templates,
		fileSystem: fileSystem,
	}
}

type djangoPlatformifier struct {
	templates  fs.FS
	fileSystem FS
}

func (p *djangoPlatformifier) Platformify(ctx context.Context, input *UserInput) error {
	appRoot := filepath.Join(input.Root, input.ApplicationRoot)
	if settingsPath := p.fileSystem.Find(appRoot, settingsPyFile, true); len(settingsPath) > 0 {
		pshSettingsPath := filepath.Join(filepath.Dir(settingsPath[0]), settingsPshPyFile)
		tpl, parseErr := template.New(settingsPshPyFile).Funcs(sprig.FuncMap()).
			ParseFS(p.templates, settingsPshPyFile)
		if parseErr != nil {
			return fmt.Errorf("could not parse template: %w", parseErr)
		}
		file, err := p.fileSystem.Create(pshSettingsPath)
		if err != nil {
			return err
		}
		defer file.Close()

		err = tpl.Execute(file, input)
		if err != nil {
			return err
		}
	}

	// append from settings_psh import * to the bottom of settings.py
	if settingsPath := p.fileSystem.Find(appRoot, settingsPyFile, true); len(settingsPath) > 0 {
		file, err := p.fileSystem.Open(settingsPath[0], os.O_APPEND|os.O_RDWR, 0o644)
		if err != nil {
			return nil
		}
		defer file.Close()

		// Check if there is an import line in the file
		found, err := utils.ContainsStringInFile(file, importSettingsPshLine)
		if err != nil {
			return err
		}

		if !found {
			if _, err = file.Write([]byte("\n\n" + importSettingsPshLine + "\n")); err != nil {
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
