package writer

import (
	"reflect"
	"testing"
	"text/template"

	"github.com/platformsh/platformify/platformifiers"
)

func TestNewWriter(t *testing.T) {
	var tests []struct {
		name    string
		want    Writer
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWriter()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWriter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWriter() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		tmplFilePath string
	}
	var tests []struct {
		name    string
		args    args
		want    *template.Template
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.tmplFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_Override(t *testing.T) {
	type args struct {
		tmpl *template.Template
	}
	var tests []struct {
		name    string
		args    args
		want    *template.Template
		wantErr bool
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := App{}
			got, err := app.Override(tt.args.tmpl)
			if (err != nil) != tt.wantErr {
				t.Errorf("Override() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Override() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_Write(t *testing.T) {
	type args struct {
		pfier platformifiers.Platformifier
	}
	var tests []struct {
		name    string
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := App{}
			if err := app.Write(tt.args.pfier); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
