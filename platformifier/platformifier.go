package platformifier

import (
	"context"
	"embed"
	"io"
	"io/fs"

	"github.com/platformsh/platformify/internal/utils"
)

var (
	//go:embed templates/**/*
	templatesFS embed.FS
)

const (
	// Path names passed to open are UTF-8-encoded,
	// unrooted, slash-separated sequences of path elements, like "x/y/z".
	// Path names MUST NOT contain an element that is "." or ".." or the empty string,
	// Paths MUST NOT start or end with a slash: "/x" and "x/" are invalid.
	genericDir = "templates/generic"
	djangoDir  = "templates/django"
	laravelDir = "templates/laravel"
	nextjsDir  = "templates/nextjs"
)

// A platformifier handles the business logic of a given runtime to platformify.
//
//go:generate mockgen -destination=mocks.go -package=platformifier -source=platformifier.go
type platformifier interface {
	// Platformify loads and writes the templates to the user's system.
	Platformify(ctx context.Context, input *UserInput) error
}

type fileCreator interface {
	Create(filePath string) (io.WriteCloser, error)
}

// New creates Platformifier with the appropriate platformifier stack based on UserInput.
func New(input *UserInput) *Platformifier {
	creator := utils.NewFileCreator()

	// fs.Sub(...) returns an error only if the given path name is invalid.
	// Since we determine the path name ourselves in advance,
	// there is no need to check for errors in this path name.
	templates, _ := fs.Sub(templatesFS, genericDir)
	stacks := []platformifier{newGenericPlatformifier(templates, creator)}

	switch input.Stack {
	case Django:
		// No need to check for errors (see the comment above)
		templates, _ = fs.Sub(templatesFS, djangoDir)
		stacks = append(stacks, newDjangoPlatformifier(templates, creator))
	case Laravel:
		// No need to check for errors (see the comment above)
		templates, _ = fs.Sub(templatesFS, laravelDir)
		stacks = append(stacks, newLaravelPlatformifier(templates, creator))
	case NextJS:
		// No need to check for errors (see the comment above)
		templates, _ = fs.Sub(templatesFS, nextjsDir)
		stacks = append(stacks, newNextJSPlatformifier(templates, creator))
	}

	return &Platformifier{
		input:  input,
		stacks: stacks,
	}
}

// A Platformifier handles the business logic of a given runtime to platformify.
type Platformifier struct {
	input  *UserInput
	stacks []platformifier
}

func (p *Platformifier) Platformify(ctx context.Context) error {
	out := make(chan error)

	go func() {
		var err error
		defer func() { out <- err }()
		for _, stack := range p.stacks {
			err = stack.Platformify(ctx, p.input)
			if err != nil {
				return
			}
		}
	}()

	// Do not wait for the end of the command execution
	// if the context has been canceled or the deadline has been exceeded
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-out:
		return err
	}
}
