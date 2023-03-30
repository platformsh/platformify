package platformifiers

import (
	"context"
	"testing"

	"github.com/platformsh/platformify/internal/models"
)

func TestNextJSPlatformifier_Platformify(t *testing.T) {
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
	}{
		{
			name:    "when the stack is empty",
			fields:  fields{answers: &models.Answers{Stack: ""}},
			wantErr: true,
		},
		{
			"when the stack is wrong",
			fields{answers: &models.Answers{Stack: "wrong"}},
			args{ctx: nil}, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			platformifier := NextJSPlatformifier{}
			p := platformifier.setUserInput(tt.fields.answers)
			if err := p.Platformify(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Platformify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
