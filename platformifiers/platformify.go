package platformifiers

import (
	"embed"
	"os"
	"path"
	"text/template"

	"github.com/platformsh/platformify/internal/models"
)

//go:embed templates/**/*
var templatesFs embed.FS

// A PlatformifierInterface describes platformifiers. A Platformifier handles the business logic of a given runtime.
type PlatformifierInterface interface {
	// setPshConfig maps answers to config values.
	setPshConfig(answers *models.Answers) *Platformifier
	// GetPshConfig is the getter for the PshConfig for the platformifier.
	GetPshConfig() PshConfig
	// getRelationships maps service names from answers to config relationships.
	getRelationships(answers *models.Answers) map[string]string
	// Platformify exports the configuration to yaml files for the user's project.
	Platformify() error
}

// GetPlatformifier is a Platformifier factory creating the appropriate instance based on UserInput.
func GetPlatformifier(answers *models.Answers) (PlatformifierInterface, error) {
	switch answers.Stack.String() {
	case models.Laravel.String():
		return NewLaravelPlatformifier(answers)
	case models.NextJS.String():
		return NewNextJSPlatformifier(answers)
	default:
		return NewPlatformifier(answers)
	}
}

func writeTemplate(tplPath string, tpl *template.Template, input any) error {
	if err := os.MkdirAll(path.Dir(tplPath), os.ModeDir|os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(tplPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return tpl.Execute(f, input)
}
