package platformifiers

import (
	"fmt"

	"github.com/platformsh/platformify/internal/models"
)

type LaravelPlatformifier struct {
	Platformifier
}

func NewLaravelPlatformifier(answers *models.Answers) (*LaravelPlatformifier, error) {
	if answers.Stack != "laravel" {
		return nil, fmt.Errorf("cannot platformify non-laravel stack: %s", answers.Stack.String())
	}

	pfier := Platformifier{}
	pfier.setPshConfig(answers)
	return &LaravelPlatformifier{pfier}, nil
}
