package platformifiers

import (
	"context"
	"fmt"
	"github.com/platformsh/platformify/internal/answer"
	"strings"
)

var stack = "general"

// AppComments are comments to add to the top of .platform.app.yaml.
type AppComments string

// Service contains the configuration for a service needed by the application.
type Service struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Disk string `json:"disk"`
}

// Mount contains the configuration for writeable directories in the app.
type Mount struct {
	Name       string
	Source     string
	SourcePath string
	Service    string
}

type PshConfig struct {
	appName string `yaml:"name"`
	appType string `yaml:"type"`
	//appSize       string            `yaml:"size"`
	relationships map[string]string `yaml:"relationships"`
	//mounts        []Mount           `yaml:"mounts"`
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

// A Platformifier handles the business logic of a given runtime to platformify.
type Platformifier struct {
	PshConfig
}

func NewPlatformifier(answers *answer.Answers) *Platformifier {
	pfier := &Platformifier{}
	pfier.setPshConfig(answers)
	return pfier
}

// setPshConfig maps answers to config values.
func (pfier Platformifier) setPshConfig(answers *answer.Answers) Platformifier {

	relationships := pfier.Relationships(answers)

	config := PshConfig{
		appName:       answers.Name,
		appType:       answers.Type.String(),
		relationships: relationships,
	}
	pfier.PshConfig = config

	return pfier

}

func (pfier Platformifier) GetPshConfig() PshConfig {
	return pfier.PshConfig
}

// Relationships returns a map of service names to their relationship names.
func (pfier Platformifier) Relationships(answers *answer.Answers) map[string]string {
	relationships := make(map[string]string)
	for _, service := range answers.Services {
		endpoint := strings.Split(service.Type.Name, ":")[0]
		relationships[service.Name] = fmt.Sprintf("%s:%s", service.Name, endpoint)
	}
	return relationships
}

func (pfier Platformifier) Platformify(context context.Context) interface{} {
	// Create template writer(s).
	// Write the files.
	return nil
}
