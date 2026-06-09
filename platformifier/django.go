package platformifier

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/platformsh/platformify/internal/utils"
	"github.com/platformsh/platformify/vendorization"
)

const (
	settingsPyFile        = "settings.py"
	settingsPshPyFile     = "settings_psh.py"
	importSettingsPshLine = "from .settings_psh import *"
)

func newDjangoPlatformifier(templates, fileSystem fs.FS) *djangoPlatformifier {
	return &djangoPlatformifier{
		templates:  templates,
		fileSystem: fileSystem,
	}
}

type djangoPlatformifier struct {
	templates  fs.FS
	fileSystem fs.FS
}

func (p *djangoPlatformifier) Platformify(ctx context.Context, input *UserInput) (map[string][]byte, error) {
	files := make(map[string][]byte)
	if settingsPath := utils.FindFile(p.fileSystem, input.ApplicationRoot, settingsPyFile); len(settingsPath) > 0 {
		pshSettingsPath := filepath.Join(filepath.Dir(settingsPath), settingsPshPyFile)
		tpl, parseErr := template.New(settingsPshPyFile).Funcs(sprig.FuncMap()).
			ParseFS(p.templates, settingsPshPyFile)
		if parseErr != nil {
			return nil, fmt.Errorf("could not parse template: %w", parseErr)
		}

		pshSettingsBuffer := &bytes.Buffer{}
		assets, _ := vendorization.FromContext(ctx)
		if err := tpl.Execute(pshSettingsBuffer, templateData{input, assets}); err != nil {
			return nil, fmt.Errorf("could not execute template: %w", parseErr)
		}
		files[pshSettingsPath] = pshSettingsBuffer.Bytes()

		settingsFile, err := fs.ReadFile(p.fileSystem, settingsPath)
		if err != nil {
			return files, nil
		}

		if !bytes.Contains(settingsFile, []byte(importSettingsPshLine)) {
			b := bytes.NewBuffer(settingsFile)
			if _, err := b.WriteString("\n\n" + importSettingsPshLine + "\n"); err == nil {
				files[settingsPath] = b.Bytes()
			}
		}
	}

	return files, nil
}
