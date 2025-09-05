package infer

import "os/exec"

func CheckHasInfer() {
	// 判断系统中，是否存在 infer 命令
	// 如果不存在，提示用户安装 infer
	_, err := exec.LookPath("infer")
	if err != nil {
		panic("Please install infer first, see https://fbinfer.com/docs/getting-started")
	}
}
