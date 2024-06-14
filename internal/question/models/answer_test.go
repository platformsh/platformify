package models

import (
	"reflect"
	"testing"

	"github.com/platformsh/platformify/platformifier"
)

func Test_getRelationships(t *testing.T) {
	type args struct {
		services []Service
	}
	tests := []struct {
		name string
		args args
		want map[string]platformifier.Relationship
	}{
		{
			name: "Simple",
			args: args{
				services: []Service{
					{
						Name: "mysql",
						Type: ServiceType{
							Name:    "mysql",
							Version: "5.7",
						},
					},
				},
			},
			want: map[string]platformifier.Relationship{
				"mysql": {Service: "mysql", Endpoint: "mysql"},
			},
		},
		{
			name: "Remapped",
			args: args{
				services: []Service{
					{
						Name: "mariadb",
						Type: ServiceType{
							Name:    "mariadb",
							Version: "14",
						},
					},
					{
						Name: "oracle-mysql",
						Type: ServiceType{
							Name:    "oracle-mysql",
							Version: "14",
						},
					},
					{
						Name: "chrome-headless",
						Type: ServiceType{
							Name:    "chrome-headless",
							Version: "14",
						},
					},
					{
						Name: "redis-persistent",
						Type: ServiceType{
							Name:    "redis-persistent",
							Version: "14",
						},
					},
				},
			},
			want: map[string]platformifier.Relationship{
				"mariadb":          {Service: "mariadb", Endpoint: "mysql"},
				"oracle-mysql":     {Service: "oracle-mysql", Endpoint: "mysql"},
				"chrome-headless":  {Service: "chrome-headless", Endpoint: "http"},
				"redis-persistent": {Service: "redis-persistent", Endpoint: "redis"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRelationships(tt.args.services); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRelationships() = %v, want %v", got, tt.want)
			}
		})
	}
}
