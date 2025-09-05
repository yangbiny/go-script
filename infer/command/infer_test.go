package infer

import (
	"testing"
)

func TestRunInferCommand(t *testing.T) {
	cmd := RunInfer()

	// 检查命令名称
	if cmd.Use != "run" {
		t.Errorf("命令名称应为 run，实际为 %s", cmd.Use)
	}

	// 检查必填参数
	args := []string{"-n", "dt-vshop", "--onlyAnalyze", "true"}
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		t.Errorf("执行命令出错: %v", err)
	}
}
