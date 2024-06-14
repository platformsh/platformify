package models

import (
	"reflect"
	"testing"
)

func Test_getRelationships(t *testing.T) {
	type args struct {
		services []Service
	}
	tests := []struct {
		name string
		args args
		want map[string]string
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
			want: map[string]string{
				"mysql": "mysql",
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
			want: map[string]string{
				"mariadb":          "mysql",
				"oracle-mysql":     "mysql",
				"chrome-headless":  "http",
				"redis-persistent": "redis",
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
