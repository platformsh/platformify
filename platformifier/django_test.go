package platformifier

import (
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

	creator   *MockfileCreator
	cwd       string
	templates fs.FS
}

func (s *PlatformifyDjangoSuiteTester) SetupTest() {
	ctrl := gomock.NewController(s.T())

	s.creator = NewMockfileCreator(ctrl)

	cwd, err := os.Getwd()
	require.NoError(s.T(), err)
	s.cwd = cwd

	templates, err := fs.Sub(testDjangoTemplatesFS, djangoTemplatesDir)
	require.NoError(s.T(), err)
	s.templates = templates
}

func (s *PlatformifyDjangoSuiteTester) TestSuccessfulFileCreation() {
	// GIVEN fake settings.py file
	settingsFilePath := path.Join(s.cwd, settingsPyFile)
	f, err := os.Create(settingsFilePath)
	require.NoError(s.T(), err)
	defer func() {
		f.Close()
		err = os.Remove(settingsFilePath)
		require.NoError(s.T(), err)
	}()
	// AND mock buffers to store PSH settings file
	buff := &MockBuffer{}
	// AND working directory is a current directory
	input := &UserInput{WorkingDirectory: s.cwd}
	// AND creation of the PSH settings file returns no errors
	s.creator.EXPECT().
		Create(gomock.Eq(path.Join(input.WorkingDirectory, settingsPshPyFile))).
		Return(buff, nil).Times(1)

	// WHEN run config files creation
	p := newDjangoPlatformifier(s.templates, s.creator)
	err = p.Platformify(context.Background(), input)
	// THEN it doesn't return any errors
	assert.NoError(s.T(), err)
	// AND the buffer contains settings file
	assert.NotEmpty(s.T(), buff)

	// WHEN check if settings file contains the line that imported psh settings file
	found, err := utils.ContainsStringInFile(settingsFilePath, importSettingsPshLine)
	// THEN it doesn't return any errors
	assert.NoError(s.T(), err)
	// AND the line is found
	assert.True(s.T(), found)
}

func (s *PlatformifyDjangoSuiteTester) TestSettingsFileNotFound() {
	// GIVEN mock buffers to store PSH settings file
	buff := &MockBuffer{}
	// AND creation of the PSH settings file returns no errors
	s.creator.EXPECT().
		Create(gomock.Eq(path.Join(s.cwd, settingsPshPyFile))).
		Return(buff, nil).Times(1)
	// AND user input is empty (because it doesn't matter if it's empty or not)
	input := &UserInput{}

	// WHEN run config files creation
	p := newDjangoPlatformifier(s.templates, s.creator)
	err := p.Platformify(context.Background(), input)
	// THEN it doesn't return any errors
	assert.NoError(s.T(), err)
	// AND the buffer is empty
	assert.Empty(s.T(), buff)
}

func (s *PlatformifyDjangoSuiteTester) TestPSHSettingsFileCreationError() {
	// GIVEN fake settings.py file
	settingsFilePath := path.Join(s.cwd, settingsPyFile)
	f, err := os.Create(settingsFilePath)
	require.NoError(s.T(), err)
	defer func() {
		f.Close()
		err = os.Remove(settingsFilePath)
		require.NoError(s.T(), err)
	}()
	// AND working directory is a current directory
	input := &UserInput{WorkingDirectory: s.cwd}
	// AND creating PSH settings file fails
	s.creator.EXPECT().
		Create(gomock.Eq(path.Join(input.WorkingDirectory, settingsPshPyFile))).
		Return(nil, errors.New("")).Times(1)

	// WHEN run config files creation
	p := newDjangoPlatformifier(s.templates, s.creator)
	err = p.Platformify(context.Background(), input)
	// THEN it fails
	assert.Error(s.T(), err)
}

func TestPlatformifyDjangoSuite(t *testing.T) {
	suite.Run(t, new(PlatformifyDjangoSuiteTester))
}
