package utils

import "testing"

func Test_openWithProjectName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "testQuery",
			args: args{
				name: "acm-uact",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			openWithProjectName(tt.args.name)
		})
	}
}
