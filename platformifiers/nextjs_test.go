package platformifiers

import (
	"testing"
)

func TestNextJSPlatformifier_Platformify(t *testing.T) {
	type fields struct {
		ui *UserInput
	}
	var tests = []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "when the stack is empty",
			fields:  fields{ui: &UserInput{Stack: ""}},
			wantErr: true,
		},
		{
			name:    "when the stack is wrong",
			fields:  fields{ui: &UserInput{Stack: "wrong"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &NextJSPlatformifier{
				UserInput: tt.fields.ui,
			}
			if err := p.Platformify(); (err != nil) != tt.wantErr {
				t.Errorf("Platformify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
