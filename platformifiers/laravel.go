package platformifiers

import (
	"fmt"
	"github.com/platformsh/platformify/internal/models/answer"
)

const laravelTemplatesPath = "templates/laravel"

type LaravelPlatformifier struct {
	Platformifier
}

func NewLaravelPlatformifier(answers *answer.Answers) (*Platformifier, error) {
	if answers.Stack != "laravel" {
		return nil, fmt.Errorf("cannot platformify non-laravel stack: %s", answers.Stack)
	}

	pfier := &Platformifier{}
	pfier.setPshConfig(answers)
	return pfier, nil
}
