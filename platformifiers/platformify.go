package platformifiers

import (
	"context"
	"embed"
	"os"
	"path"
	"strings"
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

// Service contains the configuration for a service needed by the application.
type Service struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	TypeVersions []string `json:"type_versions"`
	Disk         string   `json:"disk"`
	DiskSizes    []string `json:"disk_sizes"`
}

// UserInput contains the configuration from user input.
type UserInput struct {
	Stack             string                            `json:"stack"`
	Root              string                            `json:"root"`
	ApplicationRoot   string                            `json:"application_root"`
	Name              string                            `json:"name"`
	Type              string                            `json:"type"`
	Environment       map[string]string                 `json:"environment"`
	BuildSteps        []string                          `json:"build_steps"`
	WebCommand        string                            `json:"web_command"`
	ListenInterface   string                            `json:"listen_interface"`
	DeployCommand     []string                          `json:"deploy_command"`
	DependencyManager string                            `json:"dependency_manager"`
	Locations         map[string]map[string]interface{} `json:"locations"`
	Dependencies      map[string]map[string]string      `json:"dependencies"`
	Services          []Service
	Relationships     map[string]string
}

// A PlatformifierInterface handles the business logic of a given runtime to platformify.
type PlatformifierInterface interface {
	// getRelationships converts user answers to a Relationships map.
	getRelationships(answers *models.Answers) map[string]string
	// setUserInput converts user answers to platform.sh config values.
	setUserInput(answers *models.Answers) *Platformifier
	// getTemplatesPath is the getter for the template path on an individual platformifier.
	getTemplatesPath() string

	// Platformify loads and writes the templates to the user's system.
	Platformify(ctx context.Context) error
}

// NewPlatformifier is a Platformifier factory creating the appropriate instance based on UserInput.
func NewPlatformifier(answers *models.Answers) (PlatformifierInterface, error) {
	var p PlatformifierInterface
	switch answers.Stack {
	case models.Laravel:
		p = &LaravelPlatformifier{}
	case models.NextJS:
		p = &NextJSPlatformifier{}
	case models.Django:
		p = &DjangoPlatformifier{}
	case models.GenericStack:
		p = &GenericPlatformifier{}
	default:
		p = &Platformifier{}
	}

	return p.setUserInput(answers), nil
}

// Database returns the first service that is a database.
func (ui *UserInput) Database() string {
	for _, service := range ui.Services {
		for _, db := range databases {
			if strings.Contains(service.Type, db) {
				return service.Name
			}
		}
	}

	return ""
}

// Cache returns the first service that is a database.
func (ui *UserInput) Cache() string {
	for _, service := range ui.Services {
		for _, cache := range caches {
			if strings.Contains(service.Type, cache) {
				return service.Name
			}
		}
	}

	return ""
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
