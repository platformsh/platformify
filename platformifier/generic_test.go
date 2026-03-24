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
	genericTemplatesDir = "templates/generic"
	environmentFile     = ".environment"
	appConfigFile       = ".platform.app.yaml"
	routesConfigFile    = ".platform/routes.yaml"
	servicesConfigFile  = ".platform/services.yaml"
)

var (
	//go:embed templates/generic/*
	testGenericTemplatesFS embed.FS
)

func TestGenericPlatformifier_SuccessfulConfigsCreation(t *testing.T) {
	templates, err := fs.Sub(testGenericTemplatesFS, genericTemplatesDir)
	require.NoError(t, err)

	input := &UserInput{WorkingDirectory: fstest.MapFS{}}
	p := newGenericPlatformifier(templates, fstest.MapFS{})
	files, err := p.Platformify(context.Background(), input)
	assert.NoError(t, err)

	// The returned file map should contain the expected config files.
	assert.Contains(t, files, environmentFile)
	assert.Contains(t, files, appConfigFile)
	assert.Contains(t, files, routesConfigFile)
	assert.Contains(t, files, servicesConfigFile)

	// Each file should have non-empty content.
	assert.NotEmpty(t, files[environmentFile])
	assert.NotEmpty(t, files[appConfigFile])
	assert.NotEmpty(t, files[routesConfigFile])
	assert.NotEmpty(t, files[servicesConfigFile])
}

func TestGenericPlatformifier_EmptyInput(t *testing.T) {
	templates, err := fs.Sub(testGenericTemplatesFS, genericTemplatesDir)
	require.NoError(t, err)

	// With minimal input, Platformify should still succeed (templates render without error).
	input := &UserInput{}
	p := newGenericPlatformifier(templates, fstest.MapFS{})
	files, err := p.Platformify(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, files)
}
