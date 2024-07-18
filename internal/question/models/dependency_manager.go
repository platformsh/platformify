package models

const (
	GenericDepManager DepManager = "generic"
	Pip               DepManager = "pip"
	Poetry            DepManager = "poetry"
	Pipenv            DepManager = "pipenv"
	Composer          DepManager = "composer"
	Yarn              DepManager = "yarn"
	Npm               DepManager = "npm"
	Bundler           DepManager = "bundler"
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
	case Bundler:
		return "Bundler"
	default:
		return ""
	}
}
