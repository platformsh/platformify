/*
Package collector contains an object that can collect dependencies and combine them into one structure.

USAGE:

	c := collector.New()
	c.Add(collector.Runtime("python"))
	c.Add(collector.Stack("django"))
	c.Add(
		collector.Service("db", "mysql"),
		collector.Service("cache", "redis"),
	)
	collection := c.Collect()
*/
package collector

// Collector represents an array of dependencies.
type Collector struct {
	dependencies []Dependency
}

// New returns a new collector instance.
func New() *Collector {
	return &Collector{}
}

// Add adds the given dependencies into the internal collection.
func (c *Collector) Add(dependency ...Dependency) {
	c.dependencies = append(c.dependencies, dependency...)
}

// Collect combines all the internal dependencies into a single collection structure.
func (c *Collector) Collect() *Collection {
	collection := &Collection{}
	for _, dependencySetter := range c.dependencies {
		dependencySetter(collection)
	}

	return collection
}
