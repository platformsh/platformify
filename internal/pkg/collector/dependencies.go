package collector

// Collection describes a set of dependencies.
type Collection struct {
	Runtime  string        `json:"runtime"`
	Stack    string        `json:"stack"`
	Services []ServiceInfo `json:"services"`
}

// ServiceInfo describes a service dependency.
type ServiceInfo struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

// Dependency represents a common function for adding a new dependency.
type Dependency func(c *Collection)

// Runtime sets runtime dependency.
func Runtime(runtime string) Dependency {
	return func(c *Collection) {
		c.Runtime = runtime
	}
}

// Stack sets stack dependency.
func Stack(stack string) Dependency {
	return func(c *Collection) {
		c.Stack = stack
	}
}

// Service adds service dependency.
func Service(name, image string) Dependency {
	return func(c *Collection) {
		c.Services = append(c.Services, ServiceInfo{
			Name:  name,
			Image: image,
		})
	}
}
