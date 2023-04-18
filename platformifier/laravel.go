package platformifier

import (
	"context"
	"io/fs"
)

func newLaravelPlatformifier(templates fs.FS, file fileCreator) *laravelPlatformifier {
	return &laravelPlatformifier{
		templates: templates,
		file:      file,
	}
}

type laravelPlatformifier struct {
	templates fs.FS
	file      fileCreator
}

func (p *laravelPlatformifier) Platformify(_ context.Context, _ *UserInput) error {
	return nil
}
