package models

import (
	"reflect"
	"testing"
)

func TestRuntimeForStack(t *testing.T) {
	type args struct {
		stack Stack
	}

	tests := []struct {
		args args
		want string
	}{
		{args: args{stack: Django}, want: "python"},
		{args: args{stack: Flask}, want: "python"},
		{args: args{stack: Express}, want: "nodejs"},
		{args: args{stack: Strapi}, want: "nodejs"},
		{args: args{stack: NextJS}, want: "nodejs"},
		{args: args{stack: Laravel}, want: "php"},
	}
	for _, tt := range tests {
		t.Run(tt.args.stack.Title(), func(t *testing.T) {
			got := RuntimeForStack(tt.args.stack)
			if !reflect.DeepEqual(got.Type, tt.want) {
				t.Errorf("RuntimeForStack() = %v, want %v", got, tt.want)
			}
		})
	}
}
