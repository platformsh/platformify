package platformifier

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/platformsh/platformify/validator"
)

func TestNewPlatformifier(t *testing.T) {
	genericTemplates, err := fs.Sub(templatesFS, genericDir)
	require.NoError(t, err)
	djangoTemplates, err := fs.Sub(templatesFS, djangoDir)
	require.NoError(t, err)
	laravelTemplates, err := fs.Sub(templatesFS, laravelDir)
	require.NoError(t, err)
	nextjsTemplates, err := fs.Sub(templatesFS, nextjsDir)
	require.NoError(t, err)
	fileSystem := NewOSFileSystem("")
	tests := []struct {
		name           string
		stack          Stack
		platformifiers []platformifier
	}{
		{
			name:  "generic",
			stack: Generic,
			platformifiers: []platformifier{
				&genericPlatformifier{templates: genericTemplates, fileSystem: fileSystem},
			},
		},
		{
			name:  "django",
			stack: Django,
			platformifiers: []platformifier{
				&genericPlatformifier{templates: genericTemplates, fileSystem: fileSystem},
				&djangoPlatformifier{templates: djangoTemplates, fileSystem: fileSystem},
			},
		},
		{
			name:  "laravel",
			stack: Laravel,
			platformifiers: []platformifier{
				&genericPlatformifier{templates: genericTemplates, fileSystem: fileSystem},
				&laravelPlatformifier{templates: laravelTemplates, fileSystem: fileSystem},
			},
		},
		{
			name:  "nextjs",
			stack: NextJS,
			platformifiers: []platformifier{
				&genericPlatformifier{templates: genericTemplates, fileSystem: fileSystem},
				&nextJSPlatformifier{templates: nextjsTemplates},
			},
		},
		{
			name:  "strapi",
			stack: Strapi,
			platformifiers: []platformifier{
				&genericPlatformifier{templates: genericTemplates, fileSystem: fileSystem},
			},
		},
		{
			name:  "flask",
			stack: Flask,
			platformifiers: []platformifier{
				&genericPlatformifier{templates: genericTemplates, fileSystem: fileSystem},
			},
		},
		{
			name:  "express",
			stack: Express,
			platformifiers: []platformifier{
				&genericPlatformifier{templates: genericTemplates, fileSystem: fileSystem},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN user input with given stack
			input := &UserInput{Stack: tt.stack}

			// WHEN create new platformifier
			pfier := New(input, "platformsh")
			// THEN user input inside platformifier should be the same as given
			assert.Equal(t, input, pfier.input)
			// AND length of the platformifier's stack must be equal to the length of expected stacks
			require.Len(t, pfier.stacks, len(tt.platformifiers))
			for i := range pfier.stacks {
				// AND the type of each stack should be the same as expected
				assert.IsType(t, tt.platformifiers[i], pfier.stacks[i])
				assert.True(t, reflect.DeepEqual(tt.platformifiers[i], pfier.stacks[i]))
			}
		})
	}
}

type PlatformifySuiteTester struct {
	suite.Suite

	generic *Mockplatformifier
	django  *Mockplatformifier
	laravel *Mockplatformifier
	nextjs  *Mockplatformifier
}

func (s *PlatformifySuiteTester) SetupTest() {
	ctrl := gomock.NewController(s.T())

	s.generic = NewMockplatformifier(ctrl)
	s.django = NewMockplatformifier(ctrl)
	s.laravel = NewMockplatformifier(ctrl)
	s.nextjs = NewMockplatformifier(ctrl)
}

func (s *PlatformifySuiteTester) TestSuccessfulPlatformifying() {
	// GIVEN empty context
	ctx := context.Background()
	// AND user input is empty (because it doesn't matter if it's empty or not)
	input := &UserInput{}
	// AND platformifying generic stack returns no errors
	s.generic.EXPECT().
		Platformify(gomock.Eq(ctx), gomock.Eq(input)).
		Return(nil).AnyTimes()
	// AND platformifying django stack returns no errors
	s.django.EXPECT().
		Platformify(gomock.Eq(ctx), gomock.Eq(input)).
		Return(nil).AnyTimes()
	// AND platformifying laravel stack returns no errors
	s.laravel.EXPECT().
		Platformify(gomock.Eq(ctx), gomock.Eq(input)).
		Return(nil).AnyTimes()
	// AND platformifying nextjs stack returns no errors
	s.nextjs.EXPECT().
		Platformify(gomock.Eq(ctx), gomock.Eq(input)).
		Return(nil).AnyTimes()

	tests := []struct {
		name   string
		stacks []platformifier
	}{
		{
			name:   "empty",
			stacks: []platformifier{},
		},
		{
			name:   "generic",
			stacks: []platformifier{s.generic},
		},
		{
			name:   "django",
			stacks: []platformifier{s.django},
		},
		{
			name:   "laravel",
			stacks: []platformifier{s.laravel},
		},
		{
			name:   "nextjs",
			stacks: []platformifier{s.nextjs},
		},
		{
			name:   "generic+django",
			stacks: []platformifier{s.generic, s.django},
		},
		{
			name:   "generic+laravel",
			stacks: []platformifier{s.generic, s.laravel},
		},
		{
			name:   "generic+nextjs",
			stacks: []platformifier{s.generic, s.nextjs},
		},
		{
			name:   "all",
			stacks: []platformifier{s.generic, s.django, s.laravel, s.nextjs},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// WHEN run platformifying of the given stack
			p := Platformifier{
				input:  input,
				stacks: tt.stacks,
			}
			err := p.Platformify(ctx)
			// THEN it doesn't return any errors
			assert.NoError(s.T(), err)
		})
	}
}

func (s *PlatformifySuiteTester) TestPlatformifyingError() {
	// GIVEN empty context
	ctx := context.Background()
	// AND user input is empty (because it doesn't matter if it's empty or not)
	input := &UserInput{}
	// AND platformifying generic stack fails
	s.generic.EXPECT().
		Platformify(gomock.Eq(ctx), gomock.Eq(input)).
		Return(errors.New("")).AnyTimes()
	// AND platformifying django stack fails
	s.django.EXPECT().
		Platformify(gomock.Eq(ctx), gomock.Eq(input)).
		Return(errors.New("")).AnyTimes()
	// AND platformifying laravel stack fails
	s.laravel.EXPECT().
		Platformify(gomock.Eq(ctx), gomock.Eq(input)).
		Return(errors.New("")).AnyTimes()
	// AND platformifying nextjs stack fails
	s.nextjs.EXPECT().
		Platformify(gomock.Eq(ctx), gomock.Eq(input)).
		Return(errors.New("")).AnyTimes()

	tests := []struct {
		name   string
		stacks []platformifier
	}{
		{
			name:   "generic",
			stacks: []platformifier{s.generic},
		},
		{
			name:   "django",
			stacks: []platformifier{s.django},
		},
		{
			name:   "laravel",
			stacks: []platformifier{s.laravel},
		},
		{
			name:   "nextjs",
			stacks: []platformifier{s.nextjs},
		},
		{
			name:   "generic+django",
			stacks: []platformifier{s.generic, s.django},
		},
		{
			name:   "generic+laravel",
			stacks: []platformifier{s.generic, s.laravel},
		},
		{
			name:   "generic+nextjs",
			stacks: []platformifier{s.generic, s.nextjs},
		},
		{
			name:   "all",
			stacks: []platformifier{s.generic, s.django, s.laravel, s.nextjs},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// WHEN run platformifying of the given stack
			p := Platformifier{
				input:  input,
				stacks: tt.stacks,
			}
			err := p.Platformify(ctx)
			// THEN it fails
			assert.Error(s.T(), err)
		})
	}
}

func TestPlatformifySuite(t *testing.T) {
	suite.Run(t, new(PlatformifySuiteTester))
}

func TestPlatformifier_Platformify(t *testing.T) {
	type fields struct {
		ui *UserInput
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Django",
			fields: fields{ui: &UserInput{
				Name:  "Django",
				Type:  "python",
				Stack: Django,
			}},
		},
		{
			name: "Generic",
			fields: fields{ui: &UserInput{
				Name:  "Generic",
				Type:  "java",
				Stack: Generic,
			}},
		},
		{
			name: "Laravel",
			fields: fields{ui: &UserInput{
				Name:  "Laravel",
				Type:  "php",
				Stack: Laravel,
			}},
		},
		{
			name: "Next.js",
			fields: fields{ui: &UserInput{
				Name:  "Next.js",
				Type:  "node",
				Stack: NextJS,
			}},
		},
		{
			name: "Strapi",
			fields: fields{ui: &UserInput{
				Name:  "Strapi",
				Type:  "node",
				Stack: Strapi,
			}},
		},
		{
			name: "Flask",
			fields: fields{ui: &UserInput{
				Name:  "Flask",
				Type:  "python",
				Stack: Flask,
			}},
		},
		{
			name: "Express",
			fields: fields{ui: &UserInput{
				Name:  "Express",
				Type:  "node",
				Stack: Express,
			}},
		},
	}

	// Create a temporary directory to use as the output directory.
	tempDir, err := os.MkdirTemp("", "yaml_tests")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ctx := context.Background()
	for _, tt := range tests {
		dir, err := os.MkdirTemp(tempDir, tt.name)
		if err != nil {
			t.Fatalf("Failed to create temporary %v directory: %v", tt.name, err)
		}
		tt.fields.ui.WorkingDirectory = dir
		t.Run(tt.name, func(t *testing.T) {
			if err := New(tt.fields.ui, "platformsh").Platformify(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Platformifier.Platformify() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Validate the config.
			if err := validator.ValidateConfig(dir); (err != nil) != tt.wantErr {
				t.Errorf("Platformifier.Platformify() validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
