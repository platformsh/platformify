package models

const (
	D1024 ServiceDisk = "1024"
	D2048 ServiceDisk = "2048"
	D3072 ServiceDisk = "3072"
	D4096 ServiceDisk = "4096"
	D5120 ServiceDisk = "5120"
)

var (
	ServiceDisks = []ServiceDisk{
		D1024,
		D2048,
		D3072,
		D4096,
		D5120,
	}
)

type ServiceDisk string

func (s ServiceDisk) String() string {
	return string(s)
}

func (s ServiceDisk) Title() string {
	switch s {
	case D1024:
		return "1GB"
	case D2048:
		return "2GB"
	case D3072:
		return "3GB"
	case D4096:
		return "4GB"
	case D5120:
		return "5GB"
	default:
		return ""
	}
}
