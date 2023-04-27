package models

import (
	"fmt"
)

const (
	TCP        SocketFamily = "tcp"
	UnixSocket SocketFamily = "unix"
)

var (
	SocketFamilys = SocketFamilyList{
		TCP,
		UnixSocket,
	}
)

type SocketFamily string

func (i SocketFamily) String() string {
	return string(i)
}

func (i SocketFamily) Title() string {
	switch i {
	case TCP:
		return "TCP Port"
	case UnixSocket:
		return "Unix-socket"
	default:
		return ""
	}
}

type SocketFamilyList []SocketFamily

func (i SocketFamilyList) AllTitles() []string {
	titles := make([]string, 0, len(i))
	for _, iface := range i {
		titles = append(titles, iface.Title())
	}
	return titles
}

func (i SocketFamilyList) SocketFamilyByTitle(title string) (SocketFamily, error) {
	for _, iface := range i {
		if iface.Title() == title {
			return iface, nil
		}
	}
	return "", fmt.Errorf("listen interface by title is not found")
}
