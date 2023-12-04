package platformifier

import (
	"context"
	"embed"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	djangoTemplatesDir = "templates/django"
)

var (
	//go:embed templates/django/*
	testDjangoTemplatesFS embed.FS
)

func TestDjangoPlatformifier_SuccessfulFileCreation(t *testing.T) {
	templates, err := fs.Sub(testDjangoTemplatesFS, djangoTemplatesDir)
	require.NoError(t, err)

	// Provide a filesystem with a settings.py file.
	fileSystem := fstest.MapFS{
		"settings.py": &fstest.MapFile{
			Data: []byte("# Django settings\nDEBUG = True\n"),
		},
	}
	input := &UserInput{WorkingDirectory: fileSystem}

	p := newDjangoPlatformifier(templates, fileSystem)
	files, err := p.Platformify(context.Background(), input)
	assert.NoError(t, err)

	// The PSH settings file should have been generated.
	assert.Contains(t, files, settingsPshPyFile)
	assert.NotEmpty(t, files[settingsPshPyFile])

	// The settings.py file should contain the import line.
	assert.Contains(t, files, settingsPyFile)
	assert.Contains(t, string(files[settingsPyFile]), importSettingsPshLine)
}

func TestDjangoPlatformifier_SettingsFileNotFound(t *testing.T) {
	templates, err := fs.Sub(testDjangoTemplatesFS, djangoTemplatesDir)
	require.NoError(t, err)

	// Provide an empty filesystem (no settings.py).
	fileSystem := fstest.MapFS{}
	input := &UserInput{WorkingDirectory: fileSystem}

	p := newDjangoPlatformifier(templates, fileSystem)
	files, err := p.Platformify(context.Background(), input)
	// No error when settings.py is missing.
	assert.NoError(t, err)
	// No files should be generated.
	assert.Empty(t, files)
}

func TestDjangoPlatformifier_SettingsAlreadyImported(t *testing.T) {
	templates, err := fs.Sub(testDjangoTemplatesFS, djangoTemplatesDir)
	require.NoError(t, err)

	// Provide a filesystem where settings.py already contains the import line.
	fileSystem := fstest.MapFS{
		"settings.py": &fstest.MapFile{
			Data: []byte("# Django settings\n" + importSettingsPshLine + "\n"),
		},
	}
	input := &UserInput{WorkingDirectory: fileSystem}

	p := newDjangoPlatformifier(templates, fileSystem)
	files, err := p.Platformify(context.Background(), input)
	assert.NoError(t, err)

	// The PSH settings file should still be generated.
	assert.Contains(t, files, settingsPshPyFile)
	// The settings.py should NOT be in the output (no modification needed).
	assert.NotContains(t, files, settingsPyFile)
}
