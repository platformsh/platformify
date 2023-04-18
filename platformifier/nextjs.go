package platformifier

import (
	"context"
	"io/fs"
)

func newNextJSPlatformifier(templates fs.FS, file fileCreator) *nextJSPlatformifier {
	return &nextJSPlatformifier{
		templates: templates,
		file:      file,
	}
}

type nextJSPlatformifier struct {
	templates fs.FS
	file      fileCreator
}

func (p *nextJSPlatformifier) Platformify(_ context.Context, _ *UserInput) error {
	return nil
}
