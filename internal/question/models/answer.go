package models

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/platformsh/platformify/platformifier"
)

type Answers struct {
	Stack              Stack                             `json:"stack"`
	Flavor             string                            `json:"flavor"`
	Type               RuntimeType                       `json:"type"`
	Name               string                            `json:"name"`
	ApplicationRoot    string                            `json:"application_root"`
	Environment        map[string]string                 `json:"environment"`
	BuildSteps         []string                          `json:"build_steps"`
	WebCommand         string                            `json:"web_command"`
	SocketFamily       SocketFamily                      `json:"socket_family"`
	DeployCommand      []string                          `json:"deploy_command"`
	DependencyManagers []DepManager                      `json:"dependency_managers"`
	Dependencies       map[string]map[string]string      `json:"dependencies"`
	BuildFlavor        string                            `json:"build_flavor"`
	Disk               string                            `json:"disk"`
	Mounts             map[string]map[string]string      `json:"mounts"`
	Services           []Service                         `json:"services"`
	WorkingDirectory   string                            `json:"working_directory"`
	HasGit             bool                              `json:"has_git"`
	FilesCreated       []string                          `json:"files_created"`
	Locations          map[string]map[string]interface{} `json:"locations"`
}

type Service struct {
	Name         string        `json:"name"`
	Type         ServiceType   `json:"type"`
	TypeVersions []string      `json:"type_versions"`
	Disk         ServiceDisk   `json:"disk,omitempty"`
	DiskSizes    []ServiceDisk `json:"disk_sizes"`
}

type RuntimeType struct {
	Runtime Runtime
	Version string
}

func (t RuntimeType) String() string {
	if t.Version != "" {
		return t.Runtime.String() + ":" + t.Version
	}
	return t.Runtime.String()
}

func (t RuntimeType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

type ServiceType struct {
	Name    string
	Version string
}

func (t ServiceType) String() string {
	if t.Version != "" {
		return t.Name + ":" + t.Version
	}
	return t.Name
}

func (t ServiceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func NewAnswers() *Answers {
	return &Answers{
		Environment: make(map[string]string),
		BuildSteps:  make([]string, 0),
		Services:    make([]Service, 0),
	}
}

func (a *Answers) ToUserInput() *platformifier.UserInput {
	services := make([]platformifier.Service, 0, len(a.Services))
	for _, service := range a.Services {
		diskSizes := make([]string, 0, len(service.DiskSizes))
		for _, size := range service.DiskSizes {
			diskSizes = append(diskSizes, size.String())
		}
		services = append(services, platformifier.Service{
			Name:         service.Name,
			Type:         service.Type.String(),
			TypeVersions: service.TypeVersions,
			Disk:         service.Disk.String(),
			DiskSizes:    diskSizes,
		})
	}

	locations := map[string]map[string]interface{}{
		"/": {
			"passthru": true,
		},
	}
	for key, value := range a.Locations {
		locations[key] = value
	}

	dependencyManagers := make([]string, len(a.DependencyManagers))
	for _, dm := range a.DependencyManagers {
		dependencyManagers = append(dependencyManagers, dm.String())
	}

	return &platformifier.UserInput{
		Stack:              getStack(a.Stack),
		Root:               "",
		ApplicationRoot:    filepath.Join(string(os.PathSeparator), a.ApplicationRoot),
		Name:               a.Name,
		Type:               a.Type.String(),
		Runtime:            a.Type.Runtime.String(),
		Environment:        a.Environment,
		BuildSteps:         a.BuildSteps,
		WebCommand:         a.WebCommand,
		SocketFamily:       a.SocketFamily.String(),
		DependencyManagers: dependencyManagers,
		DeployCommand:      a.DeployCommand,
		Locations:          locations,
		Dependencies:       a.Dependencies,
		BuildFlavor:        a.BuildFlavor,
		Disk:               a.Disk,
		Mounts:             a.Mounts,
		Services:           services,
		Relationships:      getRelationships(a.Services),
		WorkingDirectory:   a.WorkingDirectory,
		HasGit:             a.HasGit,
	}
}

func getStack(answersStack Stack) platformifier.Stack {
	switch answersStack {
	case Django:
		return platformifier.Django
	case Laravel:
		return platformifier.Laravel
	case NextJS:
		return platformifier.NextJS
	case Strapi:
		return platformifier.Strapi
	case Flask:
		return platformifier.Flask
	case Express:
		return platformifier.Express
	default:
		return platformifier.Generic
	}
}

// getRelationships returns a map of service names to their relationship names.
func getRelationships(services []Service) map[string]string {
	relationships := make(map[string]string)
	for _, service := range services {
		relationships[service.Name] = ""
	}
	return relationships
}
