package platformifier

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/platformsh/platformify/internal/utils"
)

const (
	djangoTemplatesDir = "templates/django"
)

var (
	//go:embed templates/django/*
	testDjangoTemplatesFS embed.FS
)

type PlatformifyDjangoSuiteTester struct {
	suite.Suite

	cwd        string
	templates  fs.FS
	fileSystem *MockFS
}

func (s *PlatformifyDjangoSuiteTester) SetupTest() {
	ctrl := gomock.NewController(s.T())

	s.fileSystem = NewMockFS(ctrl)

	cwd, err := os.Getwd()
	require.NoError(s.T(), err)
	s.cwd = cwd

	templates, err := fs.Sub(testDjangoTemplatesFS, djangoTemplatesDir)
	require.NoError(s.T(), err)
	s.templates = templates
}

func (s *PlatformifyDjangoSuiteTester) TestSuccessfulFileCreation() {
	// GIVEN mock buffers to store settings and PSH settings files
	settingsBuff, settingsPSHBuff := &MockBuffer{}, &MockBuffer{}
	// AND working directory is a current directory
	input := &UserInput{WorkingDirectory: s.cwd}
	// AND the settings.py file exists
	s.fileSystem.EXPECT().
		Find("", settingsPyFile, true).
		Return([]string{settingsPyFile}).Times(1)
	s.fileSystem.EXPECT().
		Open(gomock.Eq(settingsPyFile), gomock.Any(), gomock.Any()).
		Return(settingsBuff, nil).Times(1)
	// AND creation of the PSH settings file returns no errors
	s.fileSystem.EXPECT().
		Create(gomock.Eq(settingsPshPyFile)).
		Return(settingsPSHBuff, nil).Times(1)

	// WHEN run config files creation
	p := newDjangoPlatformifier(s.templates, s.fileSystem)
	err := p.Platformify(context.Background(), input)
	// THEN it doesn't return any errors
	assert.NoError(s.T(), err)
	// AND the buffer contains settings file
	assert.NotEmpty(s.T(), settingsPSHBuff)

	// WHEN check if settings file contains the line that imported psh settings file
	found, err := utils.ContainsStringInFile(settingsBuff, importSettingsPshLine, false)
	// THEN it doesn't return any errors
	assert.NoError(s.T(), err)
	// AND the line is found
	assert.True(s.T(), found)
}

func (s *PlatformifyDjangoSuiteTester) TestSettingsFileNotFound() {
	// GIVEN mock buffer to store PSH settings file
	buff := &MockBuffer{}
	// AND working directory is a current directory
	input := &UserInput{WorkingDirectory: s.cwd}
	// AND the settings.py file doesn't exist
	s.fileSystem.EXPECT().
		Find("", settingsPyFile, true).
		Return([]string{}).Times(1)

	// WHEN run config files creation
	p := newDjangoPlatformifier(s.templates, s.fileSystem)
	err := p.Platformify(context.Background(), input)
	// THEN it doesn't return any errors
	assert.NoError(s.T(), err)
	// AND the buffer is empty
	assert.Empty(s.T(), buff)
}

func (s *PlatformifyDjangoSuiteTester) TestPSHSettingsFileCreationError() {
	// GIVEN working directory is a current directory
	input := &UserInput{WorkingDirectory: s.cwd}
	// AND the settings.py file exists
	s.fileSystem.EXPECT().
		Find("", settingsPyFile, true).
		Return([]string{settingsPyFile}).Times(1)
	// AND creating PSH settings file fails
	s.fileSystem.EXPECT().
		Create(gomock.Eq(settingsPshPyFile)).
		Return(nil, errors.New("")).Times(1)

	// WHEN run config files creation
	p := newDjangoPlatformifier(s.templates, s.fileSystem)
	err := p.Platformify(context.Background(), input)
	// THEN it fails
	assert.Error(s.T(), err)
}

func TestPlatformifyDjangoSuite(t *testing.T) {
	suite.Run(t, new(PlatformifyDjangoSuiteTester))
}
