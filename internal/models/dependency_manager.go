package models

import (
	"fmt"
)

const (
	GenericDepManager DepManager = "generic"
	Pip               DepManager = "pip"
	Poetry            DepManager = "poetry"
	Pipenv            DepManager = "pipenv"
	Composer          DepManager = "composer"
	Yarn              DepManager = "yarn"
	Npm               DepManager = "npm"
)

var (
	DepManagers = DepManagerList{
		Pip,
		Poetry,
		Pipenv,
		Composer,
		Yarn,
		Npm,
		GenericDepManager,
	}

	DepManagersMap = DepManagerMap{
		Python: {
			Pip, Poetry, Pipenv, GenericDepManager,
		},
		PHP: {
			Composer, GenericDepManager,
		},
		NodeJS: {
			Yarn, Npm, GenericDepManager,
		},
	}
)

type DepManager string

func (m DepManager) String() string {
	return string(m)
}

func (m DepManager) Title() string {
	switch m {
	case GenericDepManager:
		return "Other"
	case Pip:
		return "Pip"
	case Poetry:
		return "Poetry"
	case Pipenv:
		return "Pipenv"
	case Composer:
		return "Composer"
	case Yarn:
		return "Yarn"
	case Npm:
		return "Npm"
	default:
		return ""
	}
}

type DepManagerList []DepManager

func (m *DepManagerList) DepManagerByTitle(title string) (DepManager, error) {
	for i := range *m {
		if DepManagers[i].Title() == title {
			return DepManagers[i], nil
		}
	}
	return "", fmt.Errorf("dependency manager by title is not found")
}

type DepManagerMap map[Runtime][]DepManager

func (m *DepManagerMap) Titles(runtime Runtime) []string {
	managers, ok := (*m)[runtime]
	if !ok {
		return nil
	}
	titles := make([]string, 0, len(managers))
	for _, manager := range managers {
		titles = append(titles, manager.Title())
	}
	return titles
}
