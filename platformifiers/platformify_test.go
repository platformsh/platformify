package platformifiers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/platformsh/platformify/internal/models"
)

func TestGetPlatformifier(t *testing.T) {
	stacks := []struct {
		name         string
		stack        models.Stack
		pType        PlatformifierInterface
		wantErrEmpty bool
	}{
		{"Generic", models.GenericStack, new(Platformifier), false},
		{"Laravel", models.Laravel, new(LaravelPlatformifier), true},
		{"NextJS", models.NextJS, new(NextJSPlatformifier), true},
	}

	for _, stack := range stacks {
		var tests = []struct {
			name    string
			args    *models.Answers
			wantErr bool
		}{
			{"when the stack is wrong", &models.Answers{Stack: "wrong"}, true},
			{"when a platformifier is created successfully", &models.Answers{Stack: stack.stack}, false},
		}

		for _, tt := range tests {
			t.Run(fmt.Sprint(stack.name, " ", tt.name), func(t *testing.T) {
				pfier, err := GetPlatformifier(tt.args)
				if err != nil && !tt.wantErr {
					t.Errorf("GetPlatformifier error = #{err}, wantErr #{tt.wantErr}")
				}
				// Don't return a Platformifier if there's an error.
				if tt.wantErr {
					assert.Nil(t, pfier)
				} else {
					// Otherwise, make sure it's the right flavor of Platformifier.
					assert.IsType(t, stack.pType, pfier, "created object is not a %sPlatformifier", stack.name)
					// And ensure it implements the PlatformifierInterface.
					var inter interface{} = pfier
					_, pass := inter.(PlatformifierInterface)
					assert.True(t, pass, "created Platformifier but it does not implement PlatformifierInterface")
				}
			})
		}
	}
}
