package question

import (
	"context"
	"reflect"
	"testing"
	"testing/fstest"

	"github.com/platformsh/platformify/internal/question/models"
)

func TestBuildSteps_Ask(t *testing.T) {
	type args struct {
		answers models.Answers
	}
	nodeJS, _ := models.Runtimes.RuntimeByType("nodejs")
	ruby, _ := models.Runtimes.RuntimeByType("ruby")
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
				Type:               models.RuntimeType{Runtime: *nodeJS, Version: "20.0"},
				Dependencies:       map[string]map[string]string{},
				DependencyManagers: []models.DepManager{models.Yarn},
				Environment:        map[string]string{},
				WorkingDirectory:   fstest.MapFS{},
			}},
			buildSteps: []string{"yarn", "yarn exec next build"},
			wantErr:    false,
		},
		{
			name: "Next.js npm fallback",
			q:    &BuildSteps{},
			args: args{models.Answers{
				Stack:              models.NextJS,
				Type:               models.RuntimeType{Runtime: *nodeJS, Version: "20.0"},
				Dependencies:       map[string]map[string]string{},
				DependencyManagers: []models.DepManager{models.Npm},
				Environment:        map[string]string{},
				WorkingDirectory:   fstest.MapFS{},
			}},
			buildSteps: []string{"npm i", "npm exec next build"},
			wantErr:    false,
		},
		{
			name: "Bundler",
			q:    &BuildSteps{},
			args: args{models.Answers{
				Stack:              models.GenericStack,
				Type:               models.RuntimeType{Runtime: *ruby, Version: "3.3"},
				Dependencies:       map[string]map[string]string{},
				DependencyManagers: []models.DepManager{models.Bundler},
				Environment:        map[string]string{},
			}},
			buildSteps: []string{"bundle install"},
			wantErr:    false,
		},
		{
			name: "Bundler with Rails",
			q:    &BuildSteps{},
			args: args{models.Answers{
				Stack:              models.Rails,
				Type:               models.RuntimeType{Runtime: *ruby, Version: "3.3"},
				Dependencies:       map[string]map[string]string{},
				DependencyManagers: []models.DepManager{models.Bundler},
				Environment:        map[string]string{},
			}},
			buildSteps: []string{"bundle install", "bundle exec rails assets:precompile"},
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
