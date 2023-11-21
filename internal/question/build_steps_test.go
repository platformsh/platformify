package question

import (
	"context"
	"reflect"
	"testing"

	"github.com/platformsh/platformify/internal/question/models"
)

func TestBuildSteps_Ask(t *testing.T) {
	type args struct {
		answers models.Answers
	}
	tests := []struct {
		name       string
		q          *BuildSteps
		args       args
		buildSteps []string
		wantErr    bool
	}{
		{
			name: "Next.js fallback",
			q:    &BuildSteps{},
			args: args{models.Answers{
				Stack:              models.NextJS,
				Type:               models.RuntimeType{Runtime: models.NodeJS, Version: "20.0"},
				Dependencies:       map[string]map[string]string{},
				DependencyManagers: []models.DepManager{models.Yarn},
				Environment:        map[string]string{},
			}},
			buildSteps: []string{"yarn", "yarn exec next build"},
			wantErr:    false,
		},
		{
			name: "Next.js npm fallback",
			q:    &BuildSteps{},
			args: args{models.Answers{
				Stack:              models.NextJS,
				Type:               models.RuntimeType{Runtime: models.NodeJS, Version: "20.0"},
				Dependencies:       map[string]map[string]string{},
				DependencyManagers: []models.DepManager{models.Npm},
				Environment:        map[string]string{},
			}},
			buildSteps: []string{"npm i", "npm exec next build"},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &BuildSteps{}
			ctx := models.ToContext(context.Background(), &tt.args.answers)
			if err := q.Ask(ctx); (err != nil) != tt.wantErr {
				t.Errorf("BuildSteps.Ask() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.args.answers.BuildSteps, tt.buildSteps) {
				t.Errorf("BuildSteps.Ask() BuildSteps = %v, want %v", tt.args.answers.BuildSteps, tt.buildSteps)
			}
		})
	}
}
