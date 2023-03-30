package platformifiers

import (
	"context"
	"embed"
	"os"
	"path"
	"text/template"

	"github.com/platformsh/platformify/internal/models"
)

var (
	//go:embed templates/**/*
	templatesFs embed.FS
	databases   = []string{
		"mariadb",
		"mysql",
		"oracle-mysql",
		"postgresql",
	}
	caches = []string{
		"redis",
		"redis-persistent",
		"memcached",
	}
)

// A PlatformifierInterface describes platformifiers. A Platformifier handles the business logic of a given runtime.
type PlatformifierInterface interface {
	// setPshConfig maps answers to config values.
	setPshConfig(answers *models.Answers) *Platformifier
	// GetPshConfig is the getter for the PshConfig for the platformifier.
	GetPshConfig() PshConfig
	// getRelationships maps service names from answers to config relationships.
	getRelationships(answers *models.Answers) map[string]string
	// Platformify exports the configuration to yaml files for the user's project.
	Platformify(ctx context.Context) error
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

func writeTemplate(_ context.Context, tplPath string, tpl *template.Template, input any) error {
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
