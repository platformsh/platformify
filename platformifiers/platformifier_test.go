package platformifiers

import (
	"testing"

	"github.com/platformsh/platformify/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestNewPlatformifier(t *testing.T) {
	type fields struct {
		answers *models.Answers
	}
	type args struct {
	}
	var tests = []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"when the stack is empty",
			fields{&models.Answers{Stack: ""}},
			args{}, false,
		},
		{
			"when the stack is wrong",
			fields{&models.Answers{Stack: "wrong"}},
			args{}, true,
		},
		{
			"when a generic platformifier is created successfully",
			fields{&models.Answers{Stack: "generic"}},
			args{}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfier, err := NewPlatformifier(tt.fields.answers)
			if err != nil && !tt.wantErr {
				t.Errorf("NewPlatformifier error = #{err}, wantErr #{tt.wantErr}")
			}
			// Don't return a Platformifier if there's an error.
			if tt.wantErr {
				assert.Nil(t, pfier)
			} else {
				// Otherwise, make sure it's a Platformifier.
				assert.IsType(t, new(Platformifier), pfier, "created object is not a Platformifier")
				// And ensure it implements the PlatformifierInterface.
				var inter interface{} = pfier
				_, pass := inter.(PlatformifierInterface)
				assert.True(t, pass, "created Platformifier but it does not implement PlatformifierInterface")
			}
		})
	}
}
