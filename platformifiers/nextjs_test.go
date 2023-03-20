package platformifiers

import (
	"context"
	"github.com/platformsh/platformify/internal/models"
	"testing"
)

func TestNextJSPlatformifier_Platformify(t *testing.T) {
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
			p := &NextJSPlatformifier{
				UserInput: tt.fields.answers,
			}
			if err := p.Platformify(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Platformify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
