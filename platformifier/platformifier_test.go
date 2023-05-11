package platformifier

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestNewPlatformifier(t *testing.T) {
	tests := []struct {
		name           string
		stack          Stack
		platformifiers []platformifier
	}{
		{
			name:           "generic",
			stack:          Generic,
			platformifiers: []platformifier{&genericPlatformifier{}},
		},
		{
			name:           "django",
			stack:          Django,
			platformifiers: []platformifier{&genericPlatformifier{}, &djangoPlatformifier{}},
		},
		{
			name:           "laravel",
			stack:          Laravel,
			platformifiers: []platformifier{&genericPlatformifier{}, &laravelPlatformifier{}},
		},
		{
			name:           "nextjs",
			stack:          NextJS,
			platformifiers: []platformifier{&genericPlatformifier{}, &nextJSPlatformifier{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN user input with given stack
			input := &UserInput{Stack: tt.stack}

			// WHEN create new platformifier
			pfier := New(input)
			// THEN user input inside platformifier should be the same as given
			assert.Equal(t, input, pfier.input)
			// AND length of the platformifier's stack must be equal to the length of expected stacks
			require.Len(t, pfier.stacks, len(tt.platformifiers))
			for i := range pfier.stacks {
				// AND the type of each stack should be the same as expected
				assert.IsType(t, tt.platformifiers[i], pfier.stacks[i])
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
