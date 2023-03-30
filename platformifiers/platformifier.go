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

// AppComments are comments to add to the top of .platform.app.yaml.
type AppComments string

// Mount contains the configuration for writeable directories in the app.
type Mount struct {
	Name       string
	Source     string
	SourcePath string
	Service    string
}

type PshAppConfig struct {
	appName string `yaml:"name"`
	appType string `yaml:"type"`
	// appSize       string            `yaml:"size"`
	relationships map[string]string `yaml:"relationships"`
	// mounts        []Mount           `yaml:"mounts"`
	// web
	// workers
	// timezone
	// access
	// variables
	// firewall
	// build
	// dependencies
	// hooks []Hook
	// crons
	// source
	// runtime
	// additional_hosts
}

type PshConfig struct {
	AppConfig PshAppConfig
	// Routes    []models.Route
	Services []models.Service
}

// A Platformifier handles the business logic of a given runtime to platformify.
type Platformifier struct {
	PshConfig PshConfig
	Answers   *models.Answers
}

func NewPlatformifier(answers *models.Answers) (*Platformifier, error) {
	if answers.Stack.String() == "generic" || answers.Stack.String() == "" {
		pfier := &Platformifier{}
		pfier.Answers = answers
		return pfier.setPshConfig(answers), nil
	}
	return nil, fmt.Errorf("cannot platformify stack: %s", answers.Stack.String())
}

// setPshConfig maps answers to config values.
func (pfier *Platformifier) setPshConfig(answers *models.Answers) *Platformifier {
	relationships := pfier.getRelationships(answers)

	appConfig := PshAppConfig{
		appName:       answers.Name,
		appType:       answers.Type.String(),
		relationships: relationships,
	}

	pfier.PshConfig = PshConfig{
		AppConfig: appConfig,
		Services:  answers.Services,
	}

	return pfier
}

func (pfier *Platformifier) GetPshConfig() PshConfig {
	return pfier.PshConfig
}

// Relationships returns a map of service names to their relationship names.
func (pfier *Platformifier) getRelationships(answers *models.Answers) map[string]string {
	relationships := make(map[string]string)
	for _, service := range answers.Services {
		endpoint := strings.Split(service.Type.Name, ":")[0]
		relationships[service.Name] = fmt.Sprintf("%s:%s", service.Name, endpoint)
	}
	return relationships
}

func (pfier *Platformifier) Platformify(ctx context.Context) error {
	// Get working directory.
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current working directory: %w", err)
	}
	err = fs.WalkDir(templatesFs, nextjsTemplatesPath, func(filePath string, d fs.DirEntry, walkErr error) error {
		if d.IsDir() {
			return nil
		}
		tpl, er := template.New(d.Name()).Funcs(sprig.FuncMap()).ParseFS(templatesFs, filePath)
		if er != nil {
			return fmt.Errorf("could not parse template: %w", er)
		}

		filePath = path.Join(cwd, filePath[len(nextjsTemplatesPath):])
		if er := writeTemplate(ctx, filePath, tpl, pfier.PshConfig); er != nil {
			return fmt.Errorf("could not write template: %w", er)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
