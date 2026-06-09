package models

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

const (
	GenericStack Stack = iota
	Django
	Laravel
	NextJS
	Strapi
	Flask
	Express
	Rails
	Symfony
)

var (
	Stacks = StackList{
		GenericStack,
		Django,
		Laravel,
		NextJS,
		Strapi,
		Flask,
		Express,
		Rails,
		Symfony,
	}
)

type Stack int

func (s Stack) Title() string {
	switch s {
	case GenericStack:
		return "Other"
	case Django:
		return "Django"
	case Rails:
		return "Rails"
	case Laravel:
		return "Laravel"
	case NextJS:
		return "Next.js"
	case Strapi:
		return "Strapi"
	case Flask:
		return "Flask"
	case Express:
		return "Express"
	case Symfony:
		return "Symfony"
	default:
		return ""
	}
}

func (s *Stack) WriteAnswer(_ string, value interface{}) error {
	switch answer := value.(type) {
	case survey.OptionAnswer: // Select
		stack, err := Stacks.StackByTitle(answer.Value)
		if err != nil {
			return err
		}
		*s = stack
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}

type StackList []Stack

func (s StackList) AllTitles() []string {
	titles := make([]string, 0, len(s))
	for _, stack := range s {
		titles = append(titles, stack.Title())
	}
	return titles
}

func (s StackList) StackByTitle(title string) (Stack, error) {
	for _, stack := range s {
		if stack.Title() == title {
			return stack, nil
		}
	}
	return GenericStack, fmt.Errorf("stack by title is not found")
}

func RuntimeForStack(stack Stack) *Runtime {
	switch stack {
	case Django, Flask:
		if r, err := Runtimes.RuntimeByType("python"); err == nil {
			return r
		}
	case Rails:
		if r, err := Runtimes.RuntimeByType("ruby"); err == nil {
			return r
		}
		return nil
	case Laravel, Symfony:
		if r, err := Runtimes.RuntimeByType("php"); err == nil {
			return r
		}
		return nil
	case NextJS, Strapi, Express:
		if r, err := Runtimes.RuntimeByType("nodejs"); err == nil {
			return r
		}
		return nil
	}

	return nil
}
