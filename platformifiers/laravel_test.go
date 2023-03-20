package platformifiers

import (
	"context"
	"github.com/platformsh/platformify/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLaravelPlatformifier(t *testing.T) {
	type fields struct {
		answers *models.Answers
	}
	type args struct {
		ctx context.Context
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
			args{ctx: nil}, true, true,
		},
		{
			"when the stack is wrong",
			fields{&models.Answers{Stack: "wrong"}},
			args{ctx: nil}, true, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pfier, err := NewLaravelPlatformifier(tt.fields.answers)
			if err != nil && !tt.wantErr {
				t.Errorf("NewLaravelPlatformifier error = #{err}, wantErr #{tt.wantErr}")
			}
			if tt.isNil {
				assert.Nil(t, pfier)
			}
		})
	}
}
