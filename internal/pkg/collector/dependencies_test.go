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

func TestService(t *testing.T) {
	type service struct {
		name  string
		image string
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
					name:  "db",
					image: "mariadb",
				},
			},
		},
		{
			services: []service{
				{
					name:  "db",
					image: "mariadb",
				},
				{
					name:  "cache",
					image: "redis",
				},
			},
		},
		{
			services: []service{
				{
					name:  "db1",
					image: "postgresql",
				},
				{
					name:  "db2",
					image: "mariadb",
				},
				{
					name:  "cache",
					image: "redis",
				},
			},
		},
	}
	for _, tt := range tests {
		collection := &Collection{}
		for _, s := range tt.services {
			dependencySetter := Service(s.name, s.image)
			dependencySetter(collection)
		}

		if assert.Equal(t, len(tt.services), len(collection.Services)) {
			for i := range collection.Services {
				assert.Equal(t, tt.services[i].name, collection.Services[i].Name)
				assert.Equal(t, tt.services[i].image, collection.Services[i].Image)
			}
		}
	}
}
