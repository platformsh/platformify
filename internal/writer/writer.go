// Package writer does platform template writing using text/template.
package writer

import (
	"embed"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/platformsh/platformify/platformifiers"
)

type Block struct {
	name  string
	block template.Template
}

var Blocks []Block

//go:embed templates/**/*
var templatesFs embed.FS

// Writer translates data in a platformifier into the file(s) via template.
type Writer interface {
	// Override is called if a platformifier needs to use a different template block than the default.
	Override(tmpl *template.Template) (template.Template, error)
	// Write the template to the user's project.
	Write(pfier platformifiers.Platformifier) error
}

// NewWriter is a Writer factory creating a template writer based on ?? and adding template blocks.
func NewWriter() (Writer, error) {
	// @todo Decide what flavor of writer we need.
	var writer Writer

	for _, name := range blockNames {
		tmpl, parseErr := Parse("./templates/blocks/#{name}.goyaml")
		if parseErr != nil {
			log.Fatal(fmt.Errorf("could not parse %s block template: %v", name, parseErr))
		}

		Blocks = append(Blocks, Block{name, *tmpl})
	}
	return writer, nil
}

func Parse(tmplFilePath string) (*template.Template, error) {
	cwd, wdErr := os.Getwd()
	if wdErr != nil {
		log.Fatal(fmt.Errorf("could not get current working directory: %w", wdErr))
	}

	tmpl, parseErr := template.ParseFS(templatesFs, cwd+tmplFilePath)
	if parseErr != nil {
		return nil, fmt.Errorf("could not parse %s template: %w", filename, parseErr)
	}
	return tmpl, nil
}
