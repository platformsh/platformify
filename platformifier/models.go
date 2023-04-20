package platformifier

import (
	"strings"
)

var (
	databases = []string{
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

const (
	Generic Stack = iota
	Django
	Laravel
	NextJS
)

type Stack int

// UserInput contains the configuration from user input.
type UserInput struct {
	Stack             Stack
	Root              string
	ApplicationRoot   string
	Name              string
	Type              string
	Environment       map[string]string
	BuildSteps        []string
	WebCommand        string
	ListenInterface   string
	DeployCommand     []string
	DependencyManager string
	Locations         map[string]map[string]interface{}
	Dependencies      map[string]map[string]string
	BuildFlavor       string
	Services          []Service
	Relationships     map[string]string
}

// Service contains the configuration for a service needed by the application.
type Service struct {
	Name         string
	Type         string
	TypeVersions []string
	Disk         string
	DiskSizes    []string
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
