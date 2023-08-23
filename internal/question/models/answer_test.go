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
				"mysql": "mysql:mysql",
			},
		},
		{
			name: "MariaDB and Oracle",
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
				},
			},
			want: map[string]string{
				"mariadb":      "mariadb:mysql",
				"oracle-mysql": "oracle-mysql:mysql",
			},
		},
		{
			name: "Multiple",
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
						Name: "redis",
						Type: ServiceType{
							Name:    "redis",
							Version: "6",
						},
					},
				},
			},
			want: map[string]string{
				"mariadb": "mariadb:mysql",
				"redis":   "redis:redis",
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
