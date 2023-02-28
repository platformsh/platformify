package collector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuntime(t *testing.T) {
	tests := []struct {
		runtime []string
		want    string
	}{
		{runtime: []string{}, want: ""},
		{runtime: []string{"python:3.2"}, want: "python:3.2"},
		{runtime: []string{"java", "python:3.2"}, want: "python:3.2"},
		{runtime: []string{"java", "python:3.2", "golang:1.20"}, want: "golang:1.20"},
	}
	for _, tt := range tests {
		collection := &Collection{}
		for _, r := range tt.runtime {
			dependencySetter := Runtime(r)
			dependencySetter(collection)
		}
		assert.Equal(t, tt.want, collection.Runtime)
	}
}

func TestStack(t *testing.T) {
	tests := []struct {
		stack []string
		want  string
	}{
		{stack: []string{}, want: ""},
		{stack: []string{"django4"}, want: "django4"},
		{stack: []string{"nextjs", "django4"}, want: "django4"},
		{stack: []string{"nextjs", "django4", "laravel"}, want: "laravel"},
	}
	for _, tt := range tests {
		collection := &Collection{}
		for _, s := range tt.stack {
			dependencySetter := Stack(s)
			dependencySetter(collection)
		}
		assert.Equal(t, tt.want, collection.Stack)
	}
}

func TestEnvironment(t *testing.T) {
	type environment struct {
		key   string
		value string
		want  bool
	}
	tests := []struct {
		environments []environment
	}{
		{
			environments: []environment{},
		},
		{
			environments: []environment{
				{
					key:   "NAME",
					value: "value",
					want:  true,
				},
			},
		},
		{
			environments: []environment{
				{
					key:   "NAME_ONE",
					value: "value_one",
					want:  true,
				},
				{
					key:   "NAME_TWO",
					value: "value_two",
					want:  true,
				},
			},
		},
		{
			environments: []environment{
				{
					key:   "NAME_ONE",
					value: "value_one",
					want:  false,
				},
				{
					key:   "NAME_ONE",
					value: "value_two",
					want:  true,
				},
			},
		},
	}
	for _, tt := range tests {
		collection := &Collection{}
		for _, e := range tt.environments {
			dependencySetter := Environment(e.key, e.value)
			dependencySetter(collection)
		}

		for _, env := range tt.environments {
			value, ok := collection.Environment[env.key]
			assert.True(t, ok)
			if env.want {
				assert.Equal(t, env.value, value)
			} else {
				assert.NotEqual(t, env.value, value)
			}
		}
	}
}

func TestService(t *testing.T) {
	type service struct {
		name        string
		serviceType string
		disk        string
	}
	tests := []struct {
		services []service
	}{
		{
			services: []service{},
		},
		{
			services: []service{
				{
					name:        "db",
					serviceType: "mariadb",
					disk:        "1024",
				},
			},
		},
		{
			services: []service{
				{
					name:        "db",
					serviceType: "mariadb",
					disk:        "1024",
				},
				{
					name:        "cache",
					serviceType: "redis",
					disk:        "1024",
				},
			},
		},
		{
			services: []service{
				{
					name:        "db1",
					serviceType: "postgresql",
					disk:        "1024",
				},
				{
					name:        "db2",
					serviceType: "mariadb",
					disk:        "1024",
				},
				{
					name:        "cache",
					serviceType: "redis",
					disk:        "1024",
				},
			},
		},
	}
	for _, tt := range tests {
		collection := &Collection{}
		for _, s := range tt.services {
			dependencySetter := Service(s.name, s.serviceType, s.disk)
			dependencySetter(collection)
		}

		if assert.Equal(t, len(tt.services), len(collection.Services)) {
			for i := range collection.Services {
				assert.Equal(t, tt.services[i].name, collection.Services[i].Name)
				assert.Equal(t, tt.services[i].serviceType, collection.Services[i].Type)
				assert.Equal(t, tt.services[i].disk, collection.Services[i].Disk)
			}
		}
	}
}
