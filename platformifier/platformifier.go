package platformifier

import (
	"context"
	"embed"
	"io/fs"

	"github.com/platformsh/platformify/vendorization"
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
	upsunDir   = "templates/upsun"
	djangoDir  = "templates/django"
	laravelDir = "templates/laravel"
	nextjsDir  = "templates/nextjs"
)

// A platformifier handles the business logic of a given runtime to platformify.
type platformifier interface {
	// Platformify loads and returns the rendered templates that should be written to the user's system.
	Platformify(ctx context.Context, input *UserInput) (map[string][]byte, error)
}

type templateData struct {
	*UserInput
	Assets *vendorization.VendorAssets
}

// New creates Platformifier with the appropriate platformifier stack based on UserInput.
func New(input *UserInput, flavor string) *Platformifier {
	// fs.Sub(...) returns an error only if the given path name is invalid.
	// Since we determine the path name ourselves in advance,
	// there is no need to check for errors in this path name.
	stacks := []platformifier{}
	templatesDir := genericDir
	if flavor == "upsun" {
		templatesDir = upsunDir
	}

	templates, _ := fs.Sub(templatesFS, templatesDir)
	stacks = append(stacks, newGenericPlatformifier(templates, input.WorkingDirectory))

	switch input.Stack {
	case Django:
		// No need to check for errors (see the comment above)
		templates, _ := fs.Sub(templatesFS, djangoDir)
		stacks = append(stacks, newDjangoPlatformifier(templates, input.WorkingDirectory))
	case Laravel:
		// No need to check for errors (see the comment above)
		templates, _ := fs.Sub(templatesFS, laravelDir)
		stacks = append(stacks, newLaravelPlatformifier(templates, input.WorkingDirectory))
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

func (p *Platformifier) Platformify(ctx context.Context) (map[string][]byte, error) {
	files := make(map[string][]byte)
	for _, stack := range p.stacks {
		newFiles, err := stack.Platformify(ctx, p.input)
		if err != nil {
			return nil, err
		}

		for p, contents := range newFiles {
			files[p] = contents
		}
	}

	return files, nil
}
