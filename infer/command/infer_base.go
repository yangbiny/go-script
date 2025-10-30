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

func CheckHasGit() {
	// 判断系统中，是否存在 git 命令
	// 如果不存在，提示用户安装 git
	_, err := exec.LookPath("git")
	if err != nil {
		panic("Please install git first, see https://git-scm.com/book/en/v2/Getting-Started-Installing-Git")
	}
}
