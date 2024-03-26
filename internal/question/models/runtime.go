package models

import (
	"fmt"
)

const (
	DotNet Runtime = "dotnet"
	Elixir Runtime = "elixir"
	Golang Runtime = "golang"
	Java   Runtime = "java"
	Lisp   Runtime = "lisp"
	NodeJS Runtime = "nodejs"
	PHP    Runtime = "php"
	Python Runtime = "python"
	Ruby   Runtime = "ruby"
	Rust   Runtime = "rust"
)

var (
	Runtimes = RuntimeList{
		DotNet,
		Elixir,
		Golang,
		Java,
		Lisp,
		NodeJS,
		PHP,
		Python,
		Ruby,
		Rust,
	}
)

type Runtime string

func (r Runtime) String() string {
	return string(r)
}

func (r Runtime) Title() string {
	switch r {
	case DotNet:
		return "C#/.Net Core"
	case Elixir:
		return "Elixir"
	case Golang:
		return "Go"
	case Java:
		return "Java"
	case Lisp:
		return "Lisp"
	case NodeJS:
		return "JavaScript/Node.js"
	case PHP:
		return "PHP"
	case Python:
		return "Python"
	case Ruby:
		return "Ruby"
	case Rust:
		return "Rust"
	default:
		return ""
	}
}

type RuntimeList []Runtime

func (r RuntimeList) AllTitles() []string {
	titles := make([]string, 0, len(r))
	for _, runtime := range r {
		titles = append(titles, runtime.Title())
	}
	return titles
}

func (r RuntimeList) RuntimeByTitle(title string) (Runtime, error) {
	for _, runtime := range r {
		if runtime.Title() == title {
			return runtime, nil
		}
	}
	return "", fmt.Errorf("runtime by title is not found")
}
