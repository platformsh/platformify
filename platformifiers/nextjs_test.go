package platformifiers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/platformsh/platformify/internal/models"
)

func TestNewNextJSPlatformifier(t *testing.T) {
	type fields struct {
		answers *models.Answers
	}
	var tests = []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "when the stack is empty",
			fields:  fields{&models.Answers{Stack: ""}},
			wantErr: true,
		},
		{
			name:    "when the stack is wrong",
			fields:  fields{&models.Answers{Stack: "wrong"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfier, err := NewNextJSPlatformifier(tt.fields.answers)
			if err != nil && !tt.wantErr {
				t.Errorf("Platformify() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Don't return a Platformifier if there's an error.
			if tt.wantErr {
				assert.Nil(t, pfier)
			} else {
				// Otherwise, make sure it's a Platformifier.
				assert.IsType(t, new(Platformifier), pfier, "created object is not a Platformifier")
			}
		})
	}
}
