package models

import (
	"fmt"
)

const (
	HTTP       ListenInterface = "http"
	UnixSocket ListenInterface = "unix-socket"
)

var (
	ListenInterfaces = ListenInterfaceList{
		HTTP,
		UnixSocket,
	}
)

type ListenInterface string

func (i ListenInterface) String() string {
	return string(i)
}

func (i ListenInterface) Title() string {
	switch i {
	case HTTP:
		return "HTTP"
	case UnixSocket:
		return "Unix-socket"
	default:
		return ""
	}
}

type ListenInterfaceList []ListenInterface

func (i ListenInterfaceList) AllTitles() []string {
	titles := make([]string, 0, len(i))
	for _, iface := range i {
		titles = append(titles, iface.Title())
	}
	return titles
}

func (i ListenInterfaceList) ListenInterfaceByTitle(title string) (ListenInterface, error) {
	for _, iface := range i {
		if iface.Title() == title {
			return iface, nil
		}
	}
	return "", fmt.Errorf("listen interface by title is not found")
}
