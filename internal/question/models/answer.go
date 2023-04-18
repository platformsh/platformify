package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/platformsh/platformify/platformifier"
)

type Answers struct {
	Stack             Stack                        `json:"stack"`
	Type              RuntimeType                  `json:"type"`
	Name              string                       `json:"name"`
	ApplicationRoot   string                       `json:"application_root"`
	Environment       map[string]string            `json:"environment"`
	BuildSteps        []string                     `json:"build_steps"`
	WebCommand        string                       `json:"web_command"`
	ListenInterface   ListenInterface              `json:"listen_interface"`
	DeployCommand     string                       `json:"deploy_command"`
	DependencyManager DepManager                   `json:"dependency_manager"`
	Dependencies      map[string]map[string]string `json:"dependencies"`
	BuildFlavor       string                       `json:"build_flavor"`
	Services          []Service                    `json:"services"`
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

	return &platformifier.UserInput{
		Stack:             getStack(a.Stack),
		Root:              "",
		ApplicationRoot:   a.ApplicationRoot,
		Name:              a.Name,
		Type:              a.Type.String(),
		Environment:       a.Environment,
		BuildSteps:        a.BuildSteps,
		WebCommand:        a.WebCommand,
		ListenInterface:   a.ListenInterface.String(),
		DependencyManager: a.DependencyManager.String(),
		DeployCommand:     a.DeployCommand,
		Locations: map[string]map[string]interface{}{
			"/": {
				"passthru": true,
			},
		},
		Dependencies:  a.Dependencies,
		BuildFlavor:   a.BuildFlavor,
		Services:      services,
		Relationships: getRelationships(a.Services),
	}
}

func getStack(answersStack Stack) platformifier.Stack {
	var stack platformifier.Stack
	switch answersStack {
	case Django:
		stack = platformifier.Django
	case Laravel:
		stack = platformifier.Laravel
	case NextJS:
		stack = platformifier.NextJS
	default:
		stack = platformifier.Generic
	}
	return stack
}

// getRelationships returns a map of service names to their relationship names.
func getRelationships(services []Service) map[string]string {
	relationships := make(map[string]string)
	for _, service := range services {
		endpoint := strings.Split(service.Type.Name, ":")[0]
		relationships[service.Name] = fmt.Sprintf("%s:%s", service.Name, endpoint)
	}
	return relationships
}
