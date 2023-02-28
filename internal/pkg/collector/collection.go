package collector

// Collection describes a set of dependencies.
type Collection struct {
	Runtime     string            `json:"runtime"`
	Stack       string            `json:"stack"`
	Environment map[string]string `json:"environment"`
	Services    []ServiceInfo     `json:"services"`
}

// ServiceInfo describes a service dependency.
type ServiceInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Disk string `json:"disk"`
}
