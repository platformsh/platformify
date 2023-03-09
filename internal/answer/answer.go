package answer

import (
	"encoding/json"
)

type Answers struct {
	Stack             string            `json:"stack"`
	Type              PSHType           `json:"type"`
	Name              string            `json:"name"`
	ApplicationRoot   string            `json:"application_root"`
	Environment       map[string]string `json:"environment"`
	BuildSteps        []string          `json:"build_steps"`
	WebCommand        string            `json:"web_command"`
	Listen            string            `json:"listen"`
	DeployCommand     string            `json:"deploy_command"`
	DependencyManager string            `json:"dependency_manager"`
	Services          []PSHService      `json:"services"`
}

type PSHService struct {
	Name string  `json:"name"`
	Type PSHType `json:"type"`
	Disk string  `json:"disk"`
}

type PSHType struct {
	Name    string
	Version string
}

func (t PSHType) String() string {
	if t.Version != "" {
		return t.Name + ":" + t.Version
	}
	return t.Name
}

func (t PSHType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func NewAnswers() *Answers {
	return &Answers{
		Environment: make(map[string]string),
		BuildSteps:  make([]string, 0),
		Services:    make([]PSHService, 0),
	}
}
