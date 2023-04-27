package models

const (
	TCP        SocketFamily = "tcp"
	UnixSocket SocketFamily = "unix"
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
