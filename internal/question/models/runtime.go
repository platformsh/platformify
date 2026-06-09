package models

import (
	"fmt"
)

type Runtime struct {
	Name        string
	Description string
	Disk        bool
	Docs        struct {
		URL string
		Web struct {
			Commands struct {
				Start string
			}
		}
		Locations map[string]map[string]interface{}
	}
	Type     string
	Versions struct {
		Supported []string
	}
	Runtime bool
}

func (r *Runtime) String() string {
	return r.Type
}

func (r *Runtime) Title() string {
	return r.Name
}

func (r *Runtime) DefaultVersion() string {
	if len(r.Versions.Supported) > 0 {
		return r.Versions.Supported[0]
	}

	return ""
}

type RuntimeList []*Runtime

func (r RuntimeList) AllTitles() []string {
	titles := make([]string, 0, len(r))
	for _, runtime := range r {
		titles = append(titles, runtime.Title())
	}
	return titles
}

func (r RuntimeList) RuntimeByTitle(title string) (*Runtime, error) {
	for _, runtime := range r {
		if runtime.Title() == title {
			return runtime, nil
		}
	}

	return nil, fmt.Errorf("runtime by title is not found")
}

func (r RuntimeList) RuntimeByType(typ string) (*Runtime, error) {
	for _, runtime := range r {
		if runtime.Type == typ {
			return runtime, nil
		}
	}

	return nil, fmt.Errorf("runtime by type is not found")
}
