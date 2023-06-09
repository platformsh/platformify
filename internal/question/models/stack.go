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
)

var (
	Stacks = StackList{
		Django,
		Laravel,
		NextJS,
		Strapi,
		GenericStack,
	}
)

type Stack int

func (s Stack) Title() string {
	switch s {
	case GenericStack:
		return "Other"
	case Django:
		return "Django"
	case Laravel:
		return "Laravel"
	case NextJS:
		return "Next.js"
	case Strapi:
		return "Strapi"
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

func RuntimeForStack(stack Stack) Runtime {
	switch stack {
	case Django:
		return Python
	case Laravel:
		return PHP
	case NextJS, Strapi:
		return NodeJS
	default:
		return ""
	}
}
