package models

import (
	"encoding/json"
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
