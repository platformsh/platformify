package platformifiers

import (
	"fmt"
	"github.com/platformsh/platformify/internal/models"
)

const laravelTemplatesPath = "templates/laravel"

type LaravelPlatformifier struct {
	Platformifier
}

func NewLaravelPlatformifier(answers *models.Answers) (*Platformifier, error) {
	if answers.Stack != "laravel" {
		return nil, fmt.Errorf("cannot platformify non-laravel stack: %v", answers.Stack)
	}

	pfier := &Platformifier{}
	pfier.setPshConfig(answers)
	return pfier, nil
}
