package platformifiers

import (
	"context"
	"embed"
	"fmt"
	"strings"

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
	Name         string
	Type         string
	TypeVersions []string
	Disk         string
	DiskSizes    []string
}

// UserInput contains the configuration from user input.
type UserInput struct {
	Stack             string
	Root              string
	ApplicationRoot   string
	Name              string
	Type              string
	Environment       map[string]string
	BuildSteps        []string
	WebCommand        string
	SocketFamily      string
	DeployCommand     []string
	DependencyManager string
	Locations         map[string]map[string]interface{}
	Dependencies      map[string]map[string]string
	BuildFlavor       string
	Disk              string
	Mounts            map[string]map[string]string
	Services          []Service
}

// A Platformifier handles the business logic of a given runtime to platformify.
type Platformifier interface {
	Platformify(ctx context.Context) error
}

// NewPlatformifier is a Platformifier factory creating the appropriate instance based on UserInput.
func NewPlatformifier(answers *models.Answers) (Platformifier, error) {
	services := make([]Service, 0, len(answers.Services))
	for _, service := range answers.Services {
		diskSizes := make([]string, 0, len(service.DiskSizes))
		for _, size := range service.DiskSizes {
			diskSizes = append(diskSizes, size.String())
		}
		services = append(services, Service{
			Name:         service.Name,
			Type:         service.Type.String(),
			TypeVersions: service.TypeVersions,
			Disk:         service.Disk.String(),
			DiskSizes:    diskSizes,
		})
	}
	input := &UserInput{
		Stack:             answers.Stack.String(),
		Root:              "",
		ApplicationRoot:   answers.ApplicationRoot,
		Name:              answers.Name,
		Type:              answers.Type.String(),
		Environment:       answers.Environment,
		BuildSteps:        answers.BuildSteps,
		WebCommand:        answers.WebCommand,
		SocketFamily:      answers.SocketFamily.String(),
		DependencyManager: answers.DependencyManager.String(),
		DeployCommand:     answers.DeployCommand,
		Locations: map[string]map[string]interface{}{
			"/": {
				"passthru": true,
			},
		},
		Dependencies: answers.Dependencies,
		BuildFlavor:  answers.BuildFlavor,
		Disk:         answers.Disk,
		Mounts:       answers.Mounts,
		Services:     services,
	}
	switch answers.Stack {
	case models.Laravel:
		return &LaravelPlatformifier{UserInput: input}, nil
	case models.NextJS:
		return &NextJSPlatformifier{UserInput: input}, nil
	case models.Django:
		return &DjangoPlatformifier{UserInput: input}, nil
	default:
		return &GenericPlatformifier{UserInput: input}, nil
	}
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
