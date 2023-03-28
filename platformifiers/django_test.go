package platformifiers

import (
	"testing"

	"github.com/platformsh/platformify/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestNewDjangoPlatformifier(t *testing.T) {
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
			args{}, true,
		},
		{
			"when the stack is wrong",
			fields{&models.Answers{Stack: "wrong"}},
			args{}, true,
		},
		{
			"when a django platformifier is created successfully",
			fields{&models.Answers{Stack: "django"}},
			args{}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfier, err := NewDjangoPlatformifier(tt.fields.answers)
			if err != nil && !tt.wantErr {
				t.Errorf("NewDjangoPlatformifier error = #{err}, wantErr #{tt.wantErr}")
			}
			// Don't return a Platformifier if there's an error.
			if tt.wantErr {
				assert.Nil(t, pfier)
			} else {
				// Otherwise, make sure it's a LaravelPlatformifier.
				assert.IsType(t, new(DjangoPlatformifier), pfier, "created object is not a DjangoPlatformifier")
				// And ensure it implements the PlatformifierInterface.
				var inter interface{} = pfier
				_, pass := inter.(PlatformifierInterface)
				assert.True(t, pass, "created DjangoPlatformifier but it does not implement PlatformifierInterface")
			}
		})
	}
}
