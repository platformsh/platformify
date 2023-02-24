package collector

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollector_Add(t *testing.T) {
	type fields struct {
		dependencies []Dependency
	}
	type args struct {
		dependencies []Dependency
	}
	tests := []struct {
		fields fields
		args   args
	}{
		{
			fields: fields{
				dependencies: []Dependency{},
			},
			args: args{
				dependencies: []Dependency{},
			},
		},
		{
			fields: fields{
				dependencies: []Dependency{},
			},
			args: args{
				dependencies: []Dependency{Stack("django4")},
			},
		},
		{
			fields: fields{
				dependencies: []Dependency{},
			},
			args: args{
				dependencies: []Dependency{
					Runtime("python"),
					Stack("django4"),
					Service("db", "mariadb"),
				},
			},
		},
		{
			fields: fields{
				dependencies: []Dependency{
					Runtime("golang"),
				},
			},
			args: args{
				dependencies: []Dependency{},
			},
		},
		{
			fields: fields{
				dependencies: []Dependency{
					Runtime("java"),
					Stack("django4"),
					Service("db", "mariadb"),
				},
			},
			args: args{
				dependencies: []Dependency{
					Service("cache", "redis"),
					Service("proxy-db", "postgresql"),
				},
			},
		},
	}
	for _, tt := range tests {
		c := &Collector{
			dependencies: tt.fields.dependencies,
		}
		c.Add(tt.args.dependencies...)

		assert.Equal(t, len(tt.fields.dependencies)+len(tt.args.dependencies), len(c.dependencies))
	}
}

func TestCollector_Collect(t *testing.T) {
	type fields struct {
		dependencies []Dependency
	}
	tests := []struct {
		fields fields
		want   *Collection
	}{
		{
			fields: fields{
				dependencies: []Dependency{},
			},
			want: &Collection{},
		},
		{
			fields: fields{
				dependencies: []Dependency{Stack("django4")},
			},
			want: &Collection{
				Stack: "django4",
			},
		},
		{
			fields: fields{
				dependencies: []Dependency{
					Runtime("python:3.2"),
					Stack("django4"),
				},
			},
			want: &Collection{
				Runtime: "python:3.2",
				Stack:   "django4",
			},
		},
		{
			fields: fields{
				dependencies: []Dependency{
					Runtime("python:3.2"),
					Stack("django4"),
					Service("db", "postgresql:14"),
					Service("cache", "redis"),
				},
			},
			want: &Collection{
				Runtime: "python:3.2",
				Stack:   "django4",
				Services: []ServiceInfo{
					{
						Name:  "db",
						Image: "postgresql:14",
					},
					{
						Name:  "cache",
						Image: "redis",
					},
				},
			},
		},
		{
			fields: fields{
				dependencies: []Dependency{
					Runtime("python:3.2"),
					Stack("django4"),
					Stack("laravel"),
					Runtime("php"),
					Service("db", "postgresql:14"),
					Service("cache", "redis"),
					Service("proxy-db", "mariadb"),
				},
			},
			want: &Collection{
				Runtime: "php",
				Stack:   "laravel",
				Services: []ServiceInfo{
					{
						Name:  "db",
						Image: "postgresql:14",
					},
					{
						Name:  "cache",
						Image: "redis",
					},
					{
						Name:  "proxy-db",
						Image: "mariadb",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		c := &Collector{
			dependencies: tt.fields.dependencies,
		}
		got := c.Collect()
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Collect() got = %v, want %v", got, tt.want)
		}
	}
}
