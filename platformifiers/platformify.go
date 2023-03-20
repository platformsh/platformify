package platformifiers

import (
	"context"
	"embed"
	"os"
	"path"
	"text/template"

	"github.com/platformsh/platformify/internal/models"
)

//go:embed templates/**/*
var templatesFs embed.FS

// UserInput contains the configuration from user input.
type UserInput struct {
	Stack           string                            `json:"stack"`
	Root            string                            `json:"root"`
	ApplicationRoot string                            `json:"application_root"`
	Name            string                            `json:"name"`
	Type            string                            `json:"type"`
	Environment     map[string]string                 `json:"environment"`
	BuildSteps      []string                          `json:"build_steps"`
	WebCommand      string                            `json:"web_command"`
	ListenInterface string                            `json:"listen_interface"`
	DeployCommand   string                            `json:"deploy_command"`
	Locations       map[string]map[string]interface{} `json:"locations"`
	Services        []Service
}

// A PlatformifierInterface describes platformifiers. A Platformifier handles the business logic of a given runtime.
type PlatformifierInterface interface {
	// setPshConfig maps answers to config values.
	setPshConfig(answers *models.Answers) Platformifier
	// GetPshConfig is the getter for the PshConfig for the platformifier.
	getPshConfig() PshConfig
	// getRelationships maps service names from answers to config relationships.
	getRelationships(models.Answers) map[string]string
	createWriters()
	// Platformify exports the configuration to yaml files for the user's project.
	Platformify(ctx context.Context) error
}

// GetPlatformifier is a Platformifier factory creating the appropriate instance based on Answers.
func GetPlatformifier(answers *models.Answers) (*Platformifier, error) {
	switch answers.Stack {
	case "laravel":
		return NewLaravelPlatformifier(answers)
	default:
		return NewPlatformifier(answers), nil
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
