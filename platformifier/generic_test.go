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

	creator   *MockfileCreator
	cwd       string
	templates fs.FS
}

func (s *PlatformifyGenericSuiteTester) SetupTest() {
	ctrl := gomock.NewController(s.T())

	s.creator = NewMockfileCreator(ctrl)

	cwd, err := os.Getwd()
	require.NoError(s.T(), err)
	s.cwd = cwd

	templates, err := fs.Sub(testGenericTemplatesFS, genericTemplatesDir)
	require.NoError(s.T(), err)
	s.templates = templates
}

func (s *PlatformifyGenericSuiteTester) TestSuccessfulConfigsCreation() {
	// GIVEN mock buffers to store config files
	appBuff, routesBuff, servicesBuff := &MockBuffer{}, &MockBuffer{}, &MockBuffer{}
	// AND creation of the app config file returns no errors
	s.creator.EXPECT().
		Create(gomock.Eq(path.Join(s.cwd, appConfigFile))).
		Return(appBuff, nil).Times(1)
	// AND creation of the routes config file returns no errors
	s.creator.EXPECT().
		Create(gomock.Eq(path.Join(s.cwd, routesConfigFile))).
		Return(routesBuff, nil).Times(1)
	// AND creation of the services config file returns no errors
	s.creator.EXPECT().
		Create(gomock.Eq(path.Join(s.cwd, servicesConfigFile))).
		Return(servicesBuff, nil).Times(1)
	// AND user input is empty (because it doesn't matter if it's empty or not)
	input := &UserInput{}

	// WHEN run config files creation
	p := newGenericPlatformifier(s.templates, s.creator)
	err := p.Platformify(context.Background(), input)
	// THEN it doesn't return any errors
	assert.NoError(s.T(), err)
	// AND the buffers contain configs
	assert.NotEmpty(s.T(), appBuff)
	assert.NotEmpty(s.T(), routesBuff)
	assert.NotEmpty(s.T(), servicesBuff)
}

func (s *PlatformifyGenericSuiteTester) TestAppConfigCreationError() {
	// GIVEN creating app config file fails
	s.creator.EXPECT().
		Create(gomock.Eq(path.Join(s.cwd, appConfigFile))).
		Return(nil, errors.New("")).Times(1)
	// AND creating other config files work fine
	s.creator.EXPECT().
		Create(gomock.Any()).
		Return(&MockBuffer{}, nil).AnyTimes()
	// AND user input is empty (because it doesn't matter if it's empty or not)
	input := &UserInput{}

	// WHEN run config files creation
	p := newGenericPlatformifier(s.templates, s.creator)
	err := p.Platformify(context.Background(), input)
	// THEN it fails
	assert.Error(s.T(), err)
}

func (s *PlatformifyGenericSuiteTester) TestRoutesConfigCreationError() {
	// GIVEN creating routes config file fails
	s.creator.EXPECT().
		Create(gomock.Eq(path.Join(s.cwd, routesConfigFile))).
		Return(nil, errors.New("")).Times(1)
	// AND creating other config files work fine
	s.creator.EXPECT().
		Create(gomock.Any()).
		Return(&MockBuffer{}, nil).AnyTimes()
	// AND user input is empty (because it doesn't matter if it's empty or not)
	input := &UserInput{}

	// WHEN run config files creation
	p := newGenericPlatformifier(s.templates, s.creator)
	err := p.Platformify(context.Background(), input)
	// THEN it fails
	assert.Error(s.T(), err)
}

func (s *PlatformifyGenericSuiteTester) TestServicesConfigCreationError() {
	// GIVEN creating services config file fails
	s.creator.EXPECT().
		Create(gomock.Eq(path.Join(s.cwd, servicesConfigFile))).
		Return(nil, errors.New("")).Times(1)
	// AND creating other config files work fine
	s.creator.EXPECT().
		Create(gomock.Any()).
		Return(&MockBuffer{}, nil).AnyTimes()
	// AND user input is empty (because it doesn't matter if it's empty or not)
	input := &UserInput{}

	// WHEN run config files creation
	p := newGenericPlatformifier(s.templates, s.creator)
	err := p.Platformify(context.Background(), input)
	// THEN it fails
	assert.Error(s.T(), err)
}

func TestPlatformifyGenericSuite(t *testing.T) {
	suite.Run(t, new(PlatformifyGenericSuiteTester))
}
