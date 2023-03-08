package answer

type Answers struct {
	Stack             string            `json:"stack"`
	Type              string            `json:"type"`
	Name              string            `json:"name"`
	Root              string            `json:"root"`
	Environment       map[string]string `json:"environment"`
	BuildSteps        []string          `json:"build_steps"`
	WebCommand        string            `json:"web_command"`
	Listen            string            `json:"listen"`
	DeployCommand     string            `json:"deploy_command"`
	DependencyManager string            `json:"dependency_manager"`
	Services          []Service         `json:"services"`
}

type Service struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Disk string `json:"disk"`
}

func NewAnswers() *Answers {
	return &Answers{
		Environment: make(map[string]string),
		BuildSteps:  make([]string, 0),
		Services:    make([]Service, 0),
	}
}
