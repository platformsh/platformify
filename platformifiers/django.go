package platformifiers

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/platformsh/platformify/internal/colors"
	"github.com/platformsh/platformify/internal/models"
	"github.com/platformsh/platformify/internal/utils"
)

const (
	settingsPyFile        = "settings.py"
	settingsPshPyFile     = "settings_psh.py"
	importSettingsPshLine = "from settings_psh import *"
)

type DjangoPlatformifier struct {
	*UserInput
}

func (p *DjangoPlatformifier) Platformify(ctx context.Context) error {
	if p.Stack != models.Django.String() {
		return fmt.Errorf("cannot platformify non-django stack: %s", p.Stack)
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

	appRoot := path.Join(cwd, p.Root, p.ApplicationRoot)
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
	}

	// append from settings_psh import * to the bottom of settings.py
	if settingsPath := utils.FindFile(appRoot, settingsPyFile); settingsPath != "" {
		f, err := os.OpenFile(settingsPath, os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			return nil
		}
		defer f.Close()

		// Check if there is an import line in the file
		found, err := containsStringInFile(settingsPath, importSettingsPshLine)
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

func containsStringInFile(filename, target string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), target) {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}
