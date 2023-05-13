package platformifier

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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

type MockBuffer struct {
	bytes.Buffer
}

func (b *MockBuffer) Close() error {
	return nil
}

type PlatformifyGenericSuiteTester struct {
	suite.Suite

	cwd        string
	templates  fs.FS
	fileSystem *MockFS
}

func (s *PlatformifyGenericSuiteTester) SetupTest() {
	ctrl := gomock.NewController(s.T())

	cwd, err := os.Getwd()
	require.NoError(s.T(), err)
	s.cwd = cwd

	templates, err := fs.Sub(testGenericTemplatesFS, genericTemplatesDir)
	require.NoError(s.T(), err)
	s.templates = templates

	s.fileSystem = NewMockFS(ctrl)
}

func (s *PlatformifyGenericSuiteTester) TestSuccessfulConfigsCreation() {
	// GIVEN mock buffers to store config files
	envBuff, appBuff, routesBuff, servicesBuff := &MockBuffer{}, &MockBuffer{}, &MockBuffer{}, &MockBuffer{}
	// AND working directory is a current directory
	input := &UserInput{WorkingDirectory: s.cwd}
	// AND creation of the environment file returns no errors
	s.fileSystem.EXPECT().
		CreateFile(gomock.Eq(path.Join(input.WorkingDirectory, environmentFile))).
		Return(envBuff, nil).Times(1)
	// AND creation of the app config file returns no errors
	s.fileSystem.EXPECT().
		CreateFile(gomock.Eq(path.Join(input.WorkingDirectory, appConfigFile))).
		Return(appBuff, nil).Times(1)
	// AND creation of the routes config file returns no errors
	s.fileSystem.EXPECT().
		CreateFile(gomock.Eq(path.Join(input.WorkingDirectory, routesConfigFile))).
		Return(routesBuff, nil).Times(1)
	// AND creation of the services config file returns no errors
	s.fileSystem.EXPECT().
		CreateFile(gomock.Eq(path.Join(input.WorkingDirectory, servicesConfigFile))).
		Return(servicesBuff, nil).Times(1)

	// WHEN run config files creation
	p := newGenericPlatformifier(s.templates, s.fileSystem)
	err := p.Platformify(context.Background(), input)
	// THEN it doesn't return any errors
	assert.NoError(s.T(), err)
	// AND the buffers contain configs
	assert.NotEmpty(s.T(), envBuff)
	assert.NotEmpty(s.T(), appBuff)
	assert.NotEmpty(s.T(), routesBuff)
	assert.NotEmpty(s.T(), servicesBuff)
}

func (s *PlatformifyGenericSuiteTester) TestEnvironmentCreationError() {
	// GIVEN working directory is a current directory
	input := &UserInput{WorkingDirectory: s.cwd}
	// AND creating environment file fails
	s.fileSystem.EXPECT().
		CreateFile(gomock.Eq(path.Join(input.WorkingDirectory, environmentFile))).
		Return(nil, errors.New("")).Times(1)
	// AND creating other config files work fine
	s.fileSystem.EXPECT().
		CreateFile(gomock.Any()).
		Return(&MockBuffer{}, nil).AnyTimes()

	// WHEN run config files creation
	p := newGenericPlatformifier(s.templates, s.fileSystem)
	err := p.Platformify(context.Background(), input)
	// THEN it fails
	assert.Error(s.T(), err)
}

func (s *PlatformifyGenericSuiteTester) TestAppConfigCreationError() {
	// GIVEN working directory is a current directory
	input := &UserInput{WorkingDirectory: s.cwd}
	// AND creating app config file fails
	s.fileSystem.EXPECT().
		CreateFile(gomock.Eq(path.Join(input.WorkingDirectory, appConfigFile))).
		Return(nil, errors.New("")).Times(1)
	// AND creating other config files work fine
	s.fileSystem.EXPECT().
		CreateFile(gomock.Any()).
		Return(&MockBuffer{}, nil).AnyTimes()

	// WHEN run config files creation
	p := newGenericPlatformifier(s.templates, s.fileSystem)
	err := p.Platformify(context.Background(), input)
	// THEN it fails
	assert.Error(s.T(), err)
}

func (s *PlatformifyGenericSuiteTester) TestRoutesConfigCreationError() {
	// GIVEN working directory is a current directory
	input := &UserInput{WorkingDirectory: s.cwd}
	// AND creating routes config file fails
	s.fileSystem.EXPECT().
		CreateFile(gomock.Eq(path.Join(input.WorkingDirectory, routesConfigFile))).
		Return(nil, errors.New("")).Times(1)
	// AND creating other config files work fine
	s.fileSystem.EXPECT().
		CreateFile(gomock.Any()).
		Return(&MockBuffer{}, nil).AnyTimes()

	// WHEN run config files creation
	p := newGenericPlatformifier(s.templates, s.fileSystem)
	err := p.Platformify(context.Background(), input)
	// THEN it fails
	assert.Error(s.T(), err)
}

func (s *PlatformifyGenericSuiteTester) TestServicesConfigCreationError() {
	// GIVEN working directory is a current directory
	input := &UserInput{WorkingDirectory: s.cwd}
	// AND creating services config file fails
	s.fileSystem.EXPECT().
		CreateFile(gomock.Eq(path.Join(input.WorkingDirectory, servicesConfigFile))).
		Return(nil, errors.New("")).Times(1)
	// AND creating other config files work fine
	s.fileSystem.EXPECT().
		CreateFile(gomock.Any()).
		Return(&MockBuffer{}, nil).AnyTimes()

	// WHEN run config files creation
	p := newGenericPlatformifier(s.templates, s.fileSystem)
	err := p.Platformify(context.Background(), input)
	// THEN it fails
	assert.Error(s.T(), err)
}

func TestPlatformifyGenericSuite(t *testing.T) {
	suite.Run(t, new(PlatformifyGenericSuiteTester))
}
