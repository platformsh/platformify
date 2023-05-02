package utils

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Masterminds/sprig/v3"
	"golang.org/x/exp/slices"
)

var skipDirs = []string{
	"vendor",
	"node_modules",
	".next",
}

// FileExists checks if the file exists
func FileExists(searchPath, name string) bool {
	return FindFile(searchPath, name) != ""
}

// FindFile searches for the file inside the path recursively
// and returns the full path of the file if found
func FindFile(searchPath, name string) string {
	var found string
	//nolint:errcheck
	filepath.WalkDir(searchPath, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// Skip vendor directories
			if slices.Contains(skipDirs, d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if d.Name() == name {
			found = p
			return errors.New("found")
		}

		return nil
	})

	return found
}

// FindAllFiles searches for the file inside the path recursively and returns all matches
func FindAllFiles(searchPath, name string) []string {
	found := make([]string, 0)
	_ = filepath.WalkDir(searchPath, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// Skip vendor directories
			if slices.Contains(skipDirs, d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if d.Name() == name {
			found = append(found, p)
		}

		return nil
	})

	return found
}

// GetJSONKey gets a value from a JSON file, by traversing the path given
func GetJSONKey(jsonPath []string, filePath string) (value interface{}, ok bool) {
	fin, err := os.Open(filePath)
	if err != nil {
		return nil, false
	}
	defer fin.Close()

	input, err := io.ReadAll(fin)
	if err != nil {
		return nil, false
	}

	var data map[string]interface{}
	err = json.Unmarshal(input, &data)
	if err != nil {
		return nil, false
	}

	if len(jsonPath) == 0 {
		return data, true
	}

	for _, key := range jsonPath[:len(jsonPath)-1] {
		if value, ok = data[key]; !ok {
			return nil, false
		}

		if data, ok = value.(map[string]interface{}); !ok {
			return nil, false
		}
	}

	if value, ok = data[jsonPath[len(jsonPath)-1]]; !ok {
		return nil, false
	}

	return value, true
}

// WriteTemplates in the given directory, making sure the user is okay with overwriting existing files
func WriteTemplates(ctx context.Context, root string, templates map[string]*template.Template, input any) error {
	existingFiles := make([]string, 0, len(templates))
	for path := range templates {
		if st, err := os.Stat(filepath.Join(root, path)); err == nil && !st.IsDir() {
			existingFiles = append(existingFiles, path)
		}
	}

	if len(existingFiles) > 0 {
		message := "The following files already exist."
		accept := false
		for _, path := range existingFiles {
			message += "\n  * " + path
		}
		message += "\n\nDo you want to overwrite them?"
		question := &survey.Confirm{
			Message: message,
		}
		if err := survey.AskOne(question, &accept); err != nil {
			return err
		}
		if !accept {
			return fmt.Errorf("aborted by user")
		}
	}

	for path, t := range templates {
		if err := writeTemplate(ctx, filepath.Join(root, path), t, input); err != nil {
			return err
		}
	}

	return nil
}

// GatherTemplates with the given prefix inside the templates filesystem
func GatherTemplates(_ context.Context, templatesFs fs.FS, prefix string) (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template)
	err := fs.WalkDir(templatesFs, prefix, func(filePath string, d fs.DirEntry, walkErr error) error {
		if d.IsDir() {
			return nil
		}
		tpl, parseErr := template.New(d.Name()).Funcs(sprig.FuncMap()).ParseFS(templatesFs, filePath)
		if parseErr != nil {
			return fmt.Errorf("could not parse template: %w", parseErr)
		}

		filePath = filePath[len(prefix)+1:]
		templates[filePath] = tpl
		return nil
	})

	if err != nil {
		return nil, err
	}
	return templates, nil
}

// Checks if the given file contains the given string
func ContainsStringInFile(filename, target string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), target) {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

func writeTemplate(_ context.Context, tplPath string, tpl *template.Template, input any) error {
	if err := os.MkdirAll(path.Dir(tplPath), os.ModeDir|os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(tplPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return tpl.Execute(f, input)
}

func NewFileCreator() *FileCreator {
	return &FileCreator{}
}

type FileCreator struct{}

func (f *FileCreator) Create(filePath string) (io.WriteCloser, error) {
	if err := os.MkdirAll(path.Dir(filePath), os.ModeDir|os.ModePerm); err != nil {
		return nil, err
	}

	return os.Create(filePath)
}
