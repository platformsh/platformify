package platformifier

import (
	"context"
	"io/fs"
)

func newNextJSPlatformifier(templates fs.FS) *nextJSPlatformifier {
	return &nextJSPlatformifier{
		templates: templates,
	}
}

type nextJSPlatformifier struct {
	templates fs.FS
}

func (p *nextJSPlatformifier) Platformify(_ context.Context, _ *UserInput) error {
	return nil
}
