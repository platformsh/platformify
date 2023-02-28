package collector

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

// Environment adds environment variable dependency.
func Environment(key, value string) Dependency {
	return func(c *Collection) {
		if c.Environment == nil {
			c.Environment = make(map[string]string)
		}
		c.Environment[key] = value
	}
}

// Service adds service dependency.
func Service(name, serviceType, disk string) Dependency {
	return func(c *Collection) {
		c.Services = append(c.Services, ServiceInfo{
			Name: name,
			Type: serviceType,
			Disk: disk,
		})
	}
}
