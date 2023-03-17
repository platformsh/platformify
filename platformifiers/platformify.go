package platformifiers

import (
	"context"
	"fmt"
<<<<<<< HEAD
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/platformsh/platformify/internal/models"
=======
	"github.com/platformsh/platformify/internal/answer"
>>>>>>> fcb1a28 (Refactor platformifiers to not care about template writing.)
)

// A PlatformifierInterface describes platformifiers. A Platformifier handles the business logic of a given runtime.
type PlatformifierInterface interface {
	// setPshConfig maps answers to config values.
	setPshConfig(answers *answer.Answers) Platformifier
	// GetPshConfig is the getter for the PshConfig for the platformifier.
	getPshConfig() PshConfig
	// getRelationships maps service names from answers to config relationships.
	getRelationships(answers answer.Answers) map[string]string
	createWriters()
	// Platformify exports the configuration to yaml files for the user's project.
	Platformify(ctx context.Context) error
}


// GetPlatformifier is a Platformifier factory creating the appropriate instance based on UserInput.
func GetPlatformifier(answers *answer.Answers) (*Platformifier, error) {
	services := make([]Service, 0)
	for _, service := range answers.Services {
		services = append(services, Service{
			Name: service.Name,
			Type: service.Type.String(),
			Disk: service.Disk,
		})
	}

	var pfier *Platformifier
	switch answers.Stack {
	case "laravel":
		pfier, err := NewLaravelPlatformifier(answers)
		if err != nil {
			return pfier, fmt.Errorf("could not create platformifier: %s", answers.Stack)
		}
	default:
		pfier = NewPlatformifier(answers)
	}

	return pfier, nil
}
