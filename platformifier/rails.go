package platformifier

import (
	"context"
	"io/fs"
)

func newRailsPlatformifier(templates fs.FS) *railsPlatformifier {
	return &railsPlatformifier{
		templates: templates,
	}
}

type railsPlatformifier struct {
	templates fs.FS
}

func (p *railsPlatformifier) Platformify(_ context.Context, _ *UserInput) error {
	return nil
}
