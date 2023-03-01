package platformify

import (
	"context"
	"fmt"
	"github.com/platformsh/platformify/platformifiers"
)

// Service contains the configuration for a service needed by the application.
type Service struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Disk string `json:"disk"`
}

// UserInput contains the configuration from user input.
type UserInput struct {
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
	Services        []Service
}

// A Platformifier handles the business logic of a given runtime to platformify.
type Platformifier interface {
	Platformify(ctx context.Context) error
}

// NewPlatformifier is a Platformifier factory creating the appropriate instance based on UserInput.
func NewPlatformifier(input UserInput) (Platformifier, error) {
	var pfier Platformifier

	switch input.Stack {
	case "generic":
		pfier = &platformifiers.GenericPlatformifier{UserInput: input}
	case "laravel":
		pfier = &platformifiers.LaravelPlatformifier{UserInput: input}
	default:
		return nil, fmt.Errorf("cannot platformify stack: %s", input.Stack)
	}

	return pfier, nil
}
