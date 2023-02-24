package platformify

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"golang.org/x/exp/slices"
)

//go:embed templates/**/*
var genericTemplates embed.FS

var platformFiles = []string{
	"services.yaml",
	"routes.yaml",
}

// Service contains the configuration for a service needed by the application
type Service struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Disk string `json:"disk"`
}

// Platformifier contains the configuration for the application to Platformify
type Platformifier struct {
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
	Services        []Service                         `json:"services"`
}

// Platformify will generate the .platform.app.yaml and .platform/ directory
func (p *Platformifier) Platformify(ctx context.Context) error {
	if p.Stack == "generic" {
		return p.platformifyGeneric(ctx)
	}

	return fmt.Errorf("unknown stack: %s", p.Stack)
}

func (p *Platformifier) platformifyGeneric(ctx context.Context) error {
	generic := template.New("generic").Funcs(sprig.FuncMap())
	_, err := generic.ParseFS(
		genericTemplates,
		"templates/generic/.platform.app.yaml",
		"templates/generic/.platform/*.yaml",
	)
	if err != nil {
		return fmt.Errorf("could not parse generic templates: %w", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current working directory: %w", err)
	}
	for _, t := range generic.Templates() {
		tPath := path.Join(cwd, t.Name())
		if slices.Contains(platformFiles, t.Name()) {
			tPath = path.Join(cwd, ".platform", t.Name())
		}
		if err := p.writeTemplate(ctx, tPath, t); err != nil {
			return fmt.Errorf("could not write template: %w", err)
		}
	}

	return nil
}

func (p *Platformifier) writeTemplate(ctx context.Context, tPath string, t *template.Template) error {
	if err := os.MkdirAll(path.Dir(tPath), os.ModeDir|os.ModePerm); err != nil {
		return fmt.Errorf("could not create directory %s: %w", path.Dir(tPath), err)
	}

	f, err := os.Create(tPath)
	if err != nil {
		return fmt.Errorf("could not create file %s: %w", tPath, err)
	}
	defer f.Close()
	if err := t.Execute(f, p); err != nil {
		return fmt.Errorf("could not execute template %s: %w", t.Name(), err)
	}

	return nil
}

func (p *Platformifier) Relationships() map[string]string {
	relationships := make(map[string]string)
	for _, service := range p.Services {
		endpoint := strings.Split(service.Type, ":")[0]
		relationships[service.Name] = fmt.Sprintf("%s:%s", service.Name, endpoint)
	}
	return relationships
}
