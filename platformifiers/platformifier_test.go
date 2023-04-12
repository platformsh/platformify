package platformifiers

import (
	"fmt"
	"testing"

	"github.com/platformsh/platformify/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestNewPlatformifier(t *testing.T) {
	type fields struct {
		answers *models.Answers
	}

	// Test all the known stacks create a platformifier.
	for _, stack := range models.Stacks {
		t.Run(fmt.Sprintf("when the stack is %s", stack), func(t *testing.T) {
			pfier, err := NewPlatformifier(&models.Answers{Stack: stack})
			assert.Nil(t, err, "error creating platformifier")

			// Ensure it implements the PlatformifierInterface.
			var inter interface{} = pfier
			_, pass := inter.(PlatformifierInterface)
			assert.True(t, pass, "created Platformifier does not implement PlatformifierInterface")
		})
	}

	// Test the unknown stacks.
	var tests = []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"when the stack is empty",
			fields{&models.Answers{Stack: ""}}, false,
		},
		{
			"when the stack is unrecognized",
			fields{&models.Answers{Stack: "wrong"}}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfier, err := NewPlatformifier(tt.fields.answers)
			if err != nil && !tt.wantErr {
				t.Errorf("NewPlatformifier error = #{err}, wantErr #{tt.wantErr}")
			}
			// Don't return a Platformifier if there's an error.
			if tt.wantErr || err != nil {
				assert.Nil(t, pfier)
			} else {
				// Otherwise, make sure it's a Platformifier.
				assert.IsType(t, new(Platformifier), pfier, "created object is not a Platformifier")
			}
		})
	}
}
