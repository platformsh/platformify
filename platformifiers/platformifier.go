package platformifiers

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/platformsh/platformify/internal/models"
)

// A Platformifier handles the business logic of a given runtime to platformify.
type Platformifier struct {
	UserInput *UserInput
	Answers   *models.Answers
}

func (p *Platformifier) setUserInput(answers *models.Answers) *Platformifier {
	services := make([]models.Service, 0, len(answers.Services))
	for _, service := range answers.Services {
		diskSizes := service.DiskSizes
		services = append(services, models.Service{
			Name:         service.Name,
			Type:         service.Type,
			TypeVersions: service.TypeVersions,
			Disk:         service.Disk,
			DiskSizes:    diskSizes,
		})
	}
	p.UserInput = &UserInput{
		Stack:             answers.Stack.String(),
		Root:              "",
		ApplicationRoot:   answers.ApplicationRoot,
		Name:              answers.Name,
		Type:              answers.Type.String(),
		Environment:       answers.Environment,
		BuildSteps:        answers.BuildSteps,
		WebCommand:        answers.WebCommand,
		ListenInterface:   answers.ListenInterface.String(),
		DependencyManager: answers.DependencyManager.String(),
		DeployCommand:     answers.DeployCommand,
		Locations: map[string]map[string]interface{}{
			"/": {
				"passthru": true,
			},
		},
		Dependencies:  answers.Dependencies,
		Services:      services,
		Relationships: p.getRelationships(answers),
	}

	return p
}

// Relationships returns a map of service names to their relationship names.
func (p *Platformifier) getRelationships(answers *models.Answers) map[string]string {
	relationships := make(map[string]string)
	for _, service := range answers.Services {
		endpoint := strings.Split(service.Type.Name, ":")[0]
		relationships[service.Name] = fmt.Sprintf("%s:%s", service.Name, endpoint)
	}
	return relationships
}

func (p *Platformifier) Platformify(ctx context.Context) error {
	// Get working directory.
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current working directory: %w", err)
	}
	err = fs.WalkDir(templatesFs, templatesPath, func(filePath string, d fs.DirEntry, walkErr error) error {
		if d.IsDir() {
			return nil
		}
		tpl, parseErr := template.New(d.Name()).Funcs(sprig.FuncMap()).ParseFS(templatesFs, filePath)
		if parseErr != nil {
			return fmt.Errorf("could not parse template: %w", parseErr)
		}

		filePath = path.Join(cwd, filePath[len(templatesPath):])
		if writeErr := writeTemplate(ctx, filePath, tpl, p.UserInput); writeErr != nil {
			return fmt.Errorf("could not write template: %w", writeErr)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
