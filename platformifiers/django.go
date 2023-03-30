package platformifiers

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/platformsh/platformify/internal/models"
	"github.com/platformsh/platformify/internal/utils"
)

const (
	settingsPyFile    = "settings.py"
	settingsPshPyFile = "settings_psh.py"
)

type DjangoPlatformifier struct {
	Platformifier
}

func (p *DjangoPlatformifier) Platformify(ctx context.Context) error {
	if p.UserInput.Stack != models.Django.String() {
		return fmt.Errorf("cannot platformify non-django stack: %s", p.UserInput.Stack)
	}

	// Get working directory.
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current working directory: %w", err)
	}
	err = fs.WalkDir(templatesFs, "templates/django", func(filePath string, d fs.DirEntry, walkErr error) error {
		if d.IsDir() {
			return nil
		}
		tpl, parseErr := template.New(d.Name()).Funcs(sprig.FuncMap()).ParseFS(templatesFs, filePath)
		if parseErr != nil {
			return fmt.Errorf("could not parse template: %w", parseErr)
		}

		filePath = path.Join(cwd, filePath[len("templates/django"):])
		if writeErr := writeTemplate(ctx, filePath, tpl, p.UserInput); writeErr != nil {
			return fmt.Errorf("could not write template: %w", writeErr)
		}
		return nil
	})
	if err != nil {
		return err
	}

	appRoot := path.Join(cwd, p.UserInput.Root, p.UserInput.ApplicationRoot)
	if settingsPath := utils.FindFile(appRoot, settingsPyFile); settingsPath != "" {
		pshSettingsPath := filepath.Join(filepath.Dir(settingsPath), settingsPshPyFile)
		tpl, parseErr := template.New(settingsPshPyFile).Funcs(sprig.FuncMap()).ParseFS(
			templatesFs, "templates/_extras/django/settings_psh.py",
		)
		if parseErr != nil {
			return fmt.Errorf("could not parse template: %w", parseErr)
		}
		if err := writeTemplate(ctx, pshSettingsPath, tpl, p.UserInput); err != nil {
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
