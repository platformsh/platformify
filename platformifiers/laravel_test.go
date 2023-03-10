package platformifiers

import (
	"context"
	"testing"
)

func TestLaravelPlatformifier_Platformify(t *testing.T) {
	type fields struct {
		ui *UserInput
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
			"when the stack is wrong",
			fields{ui: &UserInput{Stack: "wrong"}},
			args{ctx: nil}, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &LaravelPlatformifier{
				UserInput: tt.fields.ui,
			}
			if err := p.Platformify(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Platformify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
