package platformifier

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/platformsh/platformify/validator"
)

var (
	djangoStack = &UserInput{
		Name:            "django",
		Type:            "python",
		Stack:           Django,
		Runtime:         "python-3.13",
		ApplicationRoot: "app",
		Environment: map[string]string{
			"DJANGO_SETTINGS_MODULE": "app.settings",
			"PYTHONUNBUFFERED":       "1",
		},
		BuildFlavor: "none",
		BuildSteps: []string{
			"pip install -r requirements.txt",
			"# comment here",
			"python manage.py collectstatic --noinput",
		},
		WebCommand:   "gunicorn app.wsgi",
		SocketFamily: "unix",
		DeployCommand: []string{
			"python manage.py migrate",
			"# comment here",
		},
		DependencyManagers: []string{"pip", "yarn"},
		Locations: map[string]map[string]interface{}{
			"/": {
				"passthru": true,
			},
			"/static": {
				"root":    "static",
				"expires": "1h",
				"allow":   true,
			},
		},
		Dependencies: map[string]map[string]string{
			"python": {
				"poetry": "*",
				"pip":    ">=20.0.0",
			},
			"node": {
				"yarn": "*",
				"npm":  ">=6.0.0",
			},
		},
		Disk: "1024",
		Mounts: map[string]map[string]string{
			"/.npm": {
				"source":      "local",
				"source_path": "npm",
			},
			"/.pip": {
				"source":      "local",
				"source_path": "pip",
			},
		},
		Relationships: map[string]Relationship{
			"db":    {Service: "db", Endpoint: "postgresql"},
			"mysql": {Service: "mysql", Endpoint: "mysql"},
		},
		HasGit: true,
		Services: []Service{
			{
				Name:         "db",
				Type:         "postgres",
				TypeVersions: []string{"17", "16", "15", "14", "13", "12"},
				Disk:         "1024",
				DiskSizes:    []string{"1024", "2048"},
			},
			{
				Name:         "mysql",
				Type:         "mysql",
				TypeVersions: []string{"11.0", "10.11", "10.6", "10.5", "10.4", "10.3"},
				Disk:         "1024",
				DiskSizes:    []string{"1024", "2034"},
			},
		},
	}
	genericStack = &UserInput{
		Name:  "Generic",
		Type:  "java",
		Stack: Generic,
		Environment: map[string]string{
			"JAVA": "19",
		},
		BuildFlavor: "",
		BuildSteps: []string{
			"mvn install",
		},
		WebCommand:         "tomcat",
		SocketFamily:       "tcp",
		DeployCommand:      []string{},
		DependencyManagers: []string{"mvn"},
		Locations: map[string]map[string]interface{}{
			"/": {
				"passthru": true,
			},
		},
		Dependencies: map[string]map[string]string{},
		Disk:         "1024",
		Mounts: map[string]map[string]string{
			"/.mvn": {
				"source":      "local",
				"source_path": "maven",
			},
		},
		Relationships: map[string]Relationship{
			"mysql": {Service: "mysql", Endpoint: "mysql"},
		},
		HasGit: true,
		Services: []Service{
			{
				Name:         "mysql",
				Type:         "mysql",
				TypeVersions: []string{"13", "14", "15"},
				Disk:         "1024",
				DiskSizes:    []string{"1024", "2048"},
			},
		},
	}
	laravelStack = &UserInput{
		Name:               "Laravel",
		Type:               "php",
		Stack:              Laravel,
		Runtime:            "php-8.4",
		ApplicationRoot:    "app",
		Environment:        map[string]string{},
		BuildFlavor:        "php",
		BuildSteps:         []string{},
		DeployCommand:      []string{},
		DependencyManagers: []string{"composer"},
		Locations: map[string]map[string]interface{}{
			"/": {
				"root": "index.php",
			},
		},
		Dependencies:  map[string]map[string]string{},
		Disk:          "",
		Mounts:        map[string]map[string]string{},
		Relationships: map[string]Relationship{},
		HasGit:        false,
		Services:      []Service{},
	}
	nextJSStack = &UserInput{
		Name:  "Next.js",
		Type:  "node",
		Stack: NextJS,
	}
	strapiStack = &UserInput{
		Name:  "Strapi",
		Type:  "node",
		Stack: Strapi,
	}
	flaskStack = &UserInput{
		Name:  "Flask",
		Type:  "python",
		Stack: Flask,
	}
	expressStack = &UserInput{
		Name:  "Express",
		Type:  "node",
		Stack: Express,
	}
)

func TestNewPlatformifier(t *testing.T) {
	genericTemplates, err := fs.Sub(templatesFS, genericDir)
	require.NoError(t, err)
	djangoTemplates, err := fs.Sub(templatesFS, djangoDir)
	require.NoError(t, err)
	laravelTemplates, err := fs.Sub(templatesFS, laravelDir)
	require.NoError(t, err)
	emptyFS := fstest.MapFS{}
	tests := []struct {
		name           string
		stack          Stack
		platformifiers []platformifier
	}{
		{
			name:  "generic",
			stack: Generic,
			platformifiers: []platformifier{
				&genericPlatformifier{templates: genericTemplates, fileSystem: emptyFS},
			},
		},
		{
			name:  "django",
			stack: Django,
			platformifiers: []platformifier{
				&genericPlatformifier{templates: genericTemplates, fileSystem: emptyFS},
				&djangoPlatformifier{templates: djangoTemplates, fileSystem: emptyFS},
			},
		},
		{
			name:  "laravel",
			stack: Laravel,
			platformifiers: []platformifier{
				&genericPlatformifier{templates: genericTemplates, fileSystem: emptyFS},
				&laravelPlatformifier{templates: laravelTemplates, fileSystem: emptyFS},
			},
		},
		{
			name:  "nextjs",
			stack: NextJS,
			platformifiers: []platformifier{
				&genericPlatformifier{templates: genericTemplates, fileSystem: emptyFS},
			},
		},
		{
			name:  "strapi",
			stack: Strapi,
			platformifiers: []platformifier{
				&genericPlatformifier{templates: genericTemplates, fileSystem: emptyFS},
			},
		},
		{
			name:  "flask",
			stack: Flask,
			platformifiers: []platformifier{
				&genericPlatformifier{templates: genericTemplates, fileSystem: emptyFS},
			},
		},
		{
			name:  "express",
			stack: Express,
			platformifiers: []platformifier{
				&genericPlatformifier{templates: genericTemplates, fileSystem: emptyFS},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN user input with given stack.
			input := &UserInput{Stack: tt.stack, WorkingDirectory: emptyFS}

			// WHEN create new platformifier.
			pfier := New(input, "platform")
			// THEN user input inside platformifier should be the same as given.
			assert.Equal(t, input, pfier.input)
			// AND length of the platformifier's stack must be equal to the length of expected stacks.
			require.Len(t, pfier.stacks, len(tt.platformifiers))
			for i := range pfier.stacks {
				// AND the type of each stack should be the same as expected.
				assert.IsType(t, tt.platformifiers[i], pfier.stacks[i])
			}
		})
	}
}

// mockPlatformifier is a simple test double for the platformifier interface.
type mockPlatformifier struct {
	err error
}

func (m *mockPlatformifier) Platformify(_ context.Context, _ *UserInput) (map[string][]byte, error) {
	if m.err != nil {
		return nil, m.err
	}
	return map[string][]byte{}, nil
}

func TestPlatformifier_Platformify_Success(t *testing.T) {
	ok := &mockPlatformifier{}
	tests := []struct {
		name   string
		stacks []platformifier
	}{
		{name: "empty", stacks: []platformifier{}},
		{name: "one", stacks: []platformifier{ok}},
		{name: "two", stacks: []platformifier{ok, ok}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Platformifier{
				input:  &UserInput{},
				stacks: tt.stacks,
			}
			files, err := p.Platformify(context.Background())
			assert.NoError(t, err)
			assert.NotNil(t, files)
		})
	}
}

func TestPlatformifier_Platformify_Error(t *testing.T) {
	fail := &mockPlatformifier{err: errors.New("fail")}
	ok := &mockPlatformifier{}
	tests := []struct {
		name   string
		stacks []platformifier
	}{
		{name: "single failure", stacks: []platformifier{fail}},
		{name: "first fails", stacks: []platformifier{fail, ok}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Platformifier{
				input:  &UserInput{},
				stacks: tt.stacks,
			}
			_, err := p.Platformify(context.Background())
			assert.Error(t, err)
		})
	}
}

// writeFiles writes the file map to a directory on disk.
func writeFiles(t *testing.T, dir string, files map[string][]byte) {
	t.Helper()
	for path, contents := range files {
		absPath := filepath.Join(dir, path)
		require.NoError(t, os.MkdirAll(filepath.Dir(absPath), 0o755))
		require.NoError(t, os.WriteFile(absPath, contents, 0o644))
	}
}

func TestPlatformifier_Platformify(t *testing.T) {
	type fields struct {
		ui *UserInput
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "Django", fields: fields{ui: djangoStack}},
		{name: "Generic", fields: fields{ui: genericStack}},
		{name: "Laravel", fields: fields{ui: laravelStack}},
		{name: "Next.js", fields: fields{ui: nextJSStack}},
		{name: "Strapi", fields: fields{ui: strapiStack}},
		{name: "Flask", fields: fields{ui: flaskStack}},
		{name: "Express", fields: fields{ui: expressStack}},
	}

	tempDir, err := os.MkdirTemp("", "yaml_tests")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	ctx := context.Background()
	for _, tt := range tests {
		dir, err := os.MkdirTemp(tempDir, tt.name)
		require.NoError(t, err)
		tt.fields.ui.WorkingDirectory = os.DirFS(dir)
		t.Run(tt.name, func(t *testing.T) {
			files, pErr := New(tt.fields.ui, "platform").Platformify(ctx)
			if (pErr != nil) != tt.wantErr {
				t.Errorf("Platformifier.Platformify() error = %v, wantErr %v", pErr, tt.wantErr)
			}
			writeFiles(t, dir, files)
			if err := validator.ValidateConfig(dir, "platform"); (err != nil) != tt.wantErr {
				t.Errorf("Platformifier.Platformify() validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlatformifier_Upsunify(t *testing.T) {
	type fields struct {
		ui *UserInput
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "Django", fields: fields{ui: djangoStack}},
		{name: "Generic", fields: fields{ui: genericStack}},
		{name: "Laravel", fields: fields{ui: laravelStack}},
		{name: "Next.js", fields: fields{ui: nextJSStack}},
		{name: "Strapi", fields: fields{ui: strapiStack}},
		{name: "Flask", fields: fields{ui: flaskStack}},
		{name: "Express", fields: fields{ui: expressStack}},
	}

	tempDir, err := os.MkdirTemp("", "yaml_tests")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	ctx := context.Background()
	for _, tt := range tests {
		dir, err := os.MkdirTemp(tempDir, tt.name)
		require.NoError(t, err)
		tt.fields.ui.WorkingDirectory = os.DirFS(dir)
		t.Run(tt.name, func(t *testing.T) {
			files, pErr := New(tt.fields.ui, "upsun").Platformify(ctx)
			if (pErr != nil) != tt.wantErr {
				t.Errorf("Platformifier.Platformify() error = %v, wantErr %v", pErr, tt.wantErr)
			}
			writeFiles(t, dir, files)
			if err := validator.ValidateConfig(dir, "upsun"); (err != nil) != tt.wantErr {
				t.Errorf("Platformifier.Platformify() validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
