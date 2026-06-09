package question

import (
	"context"
	"io/fs"
	"slices"
	"testing"
	"testing/fstest"

	"github.com/platformsh/platformify/discovery"
	"github.com/platformsh/platformify/internal/question/models"
)

func TestDependencyManager_Ask(t *testing.T) {
	tests := []struct {
		fileSystem fs.FS
		want       models.DepManager
	}{
		{fileSystem: fstest.MapFS{"package-lock.json": &fstest.MapFile{}}, want: models.Npm},
		{fileSystem: fstest.MapFS{"yarn.lock": &fstest.MapFile{}}, want: models.Yarn},
		{fileSystem: fstest.MapFS{"poetry.lock": &fstest.MapFile{}}, want: models.Poetry},
		{fileSystem: fstest.MapFS{"Pipfile.lock": &fstest.MapFile{}}, want: models.Pipenv},
		{fileSystem: fstest.MapFS{"composer.lock": &fstest.MapFile{}}, want: models.Composer},
		{fileSystem: fstest.MapFS{"requirements.txt": &fstest.MapFile{}}, want: models.Pip},

		// Verify sub-directories
		{fileSystem: fstest.MapFS{"sub/package-lock.json": &fstest.MapFile{}}, want: models.Npm},
		{fileSystem: fstest.MapFS{"sub/yarn.lock": &fstest.MapFile{}}, want: models.Yarn},
		{fileSystem: fstest.MapFS{"sub/poetry.lock": &fstest.MapFile{}}, want: models.Poetry},
		{fileSystem: fstest.MapFS{"sub/Pipfile.lock": &fstest.MapFile{}}, want: models.Pipenv},
		{fileSystem: fstest.MapFS{"sub/composer.lock": &fstest.MapFile{}}, want: models.Composer},
		{fileSystem: fstest.MapFS{"sub/requirements.txt": &fstest.MapFile{}}, want: models.Pip},

		// Verify multiple examples
		{fileSystem: fstest.MapFS{
			"sub/package-lock.json": &fstest.MapFile{},
			"sub/yarn.lock":         &fstest.MapFile{},
			"sub/poetry.lock":       &fstest.MapFile{},
			"sub/Pipfile.lock":      &fstest.MapFile{},
			"sub/composer.lock":     &fstest.MapFile{},
			"sub/requirements.txt":  &fstest.MapFile{},
		}, want: models.Composer},
		{fileSystem: fstest.MapFS{
			"sub/package-lock.json": &fstest.MapFile{},
			"sub/yarn.lock":         &fstest.MapFile{},
			"sub/poetry.lock":       &fstest.MapFile{},
			"sub/Pipfile.lock":      &fstest.MapFile{},
			"sub/composer.lock":     &fstest.MapFile{},
			"sub/requirements.txt":  &fstest.MapFile{},
		}, want: models.Yarn},
	}
	for _, tt := range tests {
		t.Run(tt.want.Title(), func(t *testing.T) {
			q := &DependencyManager{}
			a := models.NewAnswers()
			a.WorkingDirectory = tt.fileSystem
			a.Discoverer = discovery.New(tt.fileSystem)
			ctx := models.ToContext(context.Background(), a)

			if err := q.Ask(ctx); err != nil {
				t.Errorf("DependencyManager.Ask() error = %v", err)
			}

			if !slices.Contains(a.DependencyManagers, tt.want) {
				t.Errorf(
					"DependencyManager.Ask().DependencyManagers = %v, want %s",
					a.DependencyManagers,
					tt.want.Title(),
				)
			}
		})
	}
}
