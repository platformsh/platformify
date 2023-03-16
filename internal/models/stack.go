package models

import (
	"fmt"
)

const (
	GenericStack Stack = "generic"
	Django       Stack = "django"
	Laravel      Stack = "laravel"
	NextJS       Stack = "next-js"
)

var (
	Stacks = StackList{
		Django,
		Laravel,
		NextJS,
		GenericStack,
	}
)

type Stack string

func (s Stack) String() string {
	return string(s)
}

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
	default:
		return ""
	}
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
	return "", fmt.Errorf("stack by title is not found")
}
