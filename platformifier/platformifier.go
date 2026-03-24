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
	genericDir = "templates/generic"
	upsunDir   = "templates/upsun"
	djangoDir  = "templates/django"
	laravelDir = "templates/laravel"
)

// A platformifier handles the business logic of a given runtime to platformify.
type platformifier interface {
	// Platformify loads and returns the rendered templates.
	Platformify(ctx context.Context, input *UserInput) (map[string][]byte, error)
}

type templateData struct {
	*UserInput
	Assets *vendorization.VendorAssets
}

// New creates Platformifier with the appropriate platformifier stack based on UserInput.
func New(input *UserInput, flavor string) *Platformifier {
	stacks := []platformifier{}
	templatesDir := genericDir
	if flavor == "upsun" {
		templatesDir = upsunDir
	}

	templates, _ := fs.Sub(templatesFS, templatesDir)
	stacks = append(stacks, newGenericPlatformifier(templates, input.WorkingDirectory))

	switch input.Stack {
	case Django:
		templates, _ := fs.Sub(templatesFS, djangoDir)
		stacks = append(stacks, newDjangoPlatformifier(templates, input.WorkingDirectory))
	case Laravel:
		templates, _ := fs.Sub(templatesFS, laravelDir)
		stacks = append(stacks, newLaravelPlatformifier(templates, input.WorkingDirectory))
	}

	return &Platformifier{
		input:  input,
		stacks: stacks,
	}
}

// Platformifier handles the business logic of a given runtime to platformify.
type Platformifier struct {
	input  *UserInput
	stacks []platformifier
}

// Platformify runs all stack platformifiers and returns the collected files.
func (p *Platformifier) Platformify(ctx context.Context) (map[string][]byte, error) {
	files := make(map[string][]byte)
	for _, stack := range p.stacks {
		newFiles, err := stack.Platformify(ctx, p.input)
		if err != nil {
			return nil, err
		}
		for path, contents := range newFiles {
			files[path] = contents
		}
	}
	return files, nil
}
