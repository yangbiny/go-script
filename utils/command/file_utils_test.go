package utils

import "testing"

func Test_csvToExcel(t *testing.T) {
	type args struct {
		csvReadPath    string
		excelWritePath string
		separator      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "testCsvToExcel",
			args: args{
				csvReadPath:    "/Users/reasonknow/Desktop/test.csv",
				excelWritePath: "/Users/reasonknow/Desktop/text.xlsx",
				separator:      ",",
			},
			want: "转换成功",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := csvToExcel(tt.args.csvReadPath, tt.args.excelWritePath, tt.args.separator); got != tt.want {
				t.Errorf("csvToExcel() = %v, want %v", got, tt.want)
			}
		})
	}
}
