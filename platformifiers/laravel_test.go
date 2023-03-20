package platformifiers

import (
	"testing"

	"github.com/platformsh/platformify/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestLaravelPlatformifier(t *testing.T) {
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
		isNil   bool
	}{
		{
			"when the stack is empty",
			fields{&models.Answers{Stack: ""}},
			args{}, true, true,
		},
		{
			"when the stack is wrong",
			fields{&models.Answers{Stack: "wrong"}},
			args{}, true, true,
		},
		{
			"when a laravel platformifier is created successfully",
			fields{&models.Answers{Stack: "laravel"}},
			args{}, false, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfier, err := NewLaravelPlatformifier(tt.fields.answers)
			if err != nil && !tt.wantErr {
				t.Errorf("NewLaravelPlatformifier error = #{err}, wantErr #{tt.wantErr}")
			}
			// Don't return a Platformifier if there's an error.
			if tt.isNil {
				assert.Nil(t, pfier)
			} else {
				assert.IsType(t, new(Platformifier), pfier, "created object is not a Platformifier")
			}
		})
	}
}
