package command

import "testing"

func Test_createChatCompletion(t *testing.T) {
	type args struct {
		configPath string
		content    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				configPath: "/Users/reasonknow/conf/chatgpt/chatgpt.conf",
				content:    "今天有这么热吗",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createChatCompletion(tt.args.configPath, tt.args.content)
		})
	}
}
