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

	"github.com/platformsh/platformify/internal/utils"
)

const (
	settingsPyFile    = "settings.py"
	settingsPshPyFile = "settings_psh.py"
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

func (p *djangoPlatformifier) Platformify(_ context.Context, input *UserInput) error {
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

		fmt.Printf(
			"We have created a %s file for you. Please add the following line to your %s file:\n",
			settingsPshPyFile,
			settingsPyFile,
		)
		fmt.Println("    from .settings_psh import *")
	}

	return nil
}
