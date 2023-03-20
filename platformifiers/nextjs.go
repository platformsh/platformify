package platformifiers

import (
	"fmt"
	"github.com/platformsh/platformify/internal/models"
)

const nextjsTemplatesPath = "templates/nextjs"

type NextJSPlatformifier struct {
	Platformifier
}

func NewNextJSPlatformifier(answers *models.Answers) (*Platformifier, error) {
	if answers.Stack.String() != models.NextJS.String() {
		return nil, fmt.Errorf("cannot platformify non-next.js stack: %v", answers.Stack)
	}
	pfier := &Platformifier{}
	pfier.setPshConfig(answers)
	return pfier, nil
}
