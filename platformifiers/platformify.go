package platformifiers

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/platformsh/platformify/internal/models"
)

//go:embed templates/**/*
var templatesFs embed.FS

// Service contains the configuration for a service needed by the application.
type Service struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Disk string `json:"disk"`
}

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

// A Platformifier handles the business logic of a given runtime to platformify.
type Platformifier interface {
	Platformify(ctx context.Context) error
}

// NewPlatformifier is a Platformifier factory creating the appropriate instance based on UserInput.
func NewPlatformifier(answers *models.Answers) (Platformifier, error) {
	services := make([]Service, 0)
	for _, service := range answers.Services {
		services = append(services, Service{
			Name: service.Name,
			Type: service.Type.String(),
			Disk: service.Disk,
		})
	}
	input := &UserInput{
		Stack:           answers.Stack.String(),
		Root:            "",
		ApplicationRoot: answers.ApplicationRoot,
		Name:            answers.Name,
		Type:            answers.Type.String(),
		Environment:     answers.Environment,
		BuildSteps:      answers.BuildSteps,
		WebCommand:      answers.WebCommand,
		ListenInterface: answers.ListenInterface,
		DeployCommand:   answers.DeployCommand,
		Locations: map[string]map[string]interface{}{
			"/": {
				"passthrough": true,
			},
		},
		Services: services,
	}
	var pfier Platformifier
	switch input.Stack {
	case "laravel":
		pfier = &LaravelPlatformifier{UserInput: input}
	case "next.js":
		pfier = &NextJSPlatformifier{UserInput: input}
	default:
		pfier = &GenericPlatformifier{UserInput: input}
	}

	return pfier, nil
}

// Relationships returns a map of service names to their relationship names.
func (ui *UserInput) Relationships() map[string]string {
	relationships := make(map[string]string)
	for _, service := range ui.Services {
		endpoint := strings.Split(service.Type, ":")[0]
		relationships[service.Name] = fmt.Sprintf("%s:%s", service.Name, endpoint)
	}
	return relationships
}

func writeTemplate(ctx context.Context, tplPath string, tpl *template.Template, input any) error {
	if err := os.MkdirAll(path.Dir(tplPath), os.ModeDir|os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(tplPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := tpl.Execute(f, input); err != nil {
		return err
	}

	return nil
}
