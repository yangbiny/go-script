package infer

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func RunInfer() *cobra.Command {
	var projectName string
	var projectPath string
	var onlyAnalyze bool

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run facebook infer on the given project",
		Args: func(cmd *cobra.Command, args []string) error {
			CheckHasInfer()
			if len(projectName) <= 0 {
				return errors.New("please provide the project name")
			}
			if len(projectPath) <= 0 {
				projectPath = "/Users/knowreason/object"
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			runInfer(projectPath, projectName, onlyAnalyze, false)
		},
	}

	cmd.Flags().StringVarP(&projectName, "projectName", "n", "", "the project name")
	cmd.Flags().StringVarP(&projectPath, "projectPath", "p", "", "the path to the project")
	cmd.Flags().BoolVarP(&onlyAnalyze, "onlyAnalyze", "", false, "only run analyze step")
	if err := cmd.MarkFlagRequired("projectName"); err != nil {
		panic(err)
	}
	return cmd
}

func RunDiffInfer() *cobra.Command {
	var onlyAnalyze bool

	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Run facebook infer on the given project with diff analysis",
		Args: func(cmd *cobra.Command, args []string) error {
			CheckHasGit()
			CheckHasInfer()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			runDiffInfer(onlyAnalyze)
		},
	}
	cmd.Flags().BoolVarP(&onlyAnalyze, "onlyAnalyze", "", false, "only run analyze step")
	return cmd
}

func runDiffInfer(onlyAnalyze bool) {
	// 存储结构：
	//1.  /Users/knowreason/soft/infer/out/diff： 存放diff 的文件类，文件名称为 项目名称_diff.txt
	//2.  /Users/knowreason/soft/infer/out/projectName： 存放 infer capture 的结果

	//1. git diff，获得改动的文件信息
	//1.1 主动探测当前的git 分支
	gitBranchCmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	branchInfo, err := gitBranchCmd.Output()
	if err != nil {
		panic("Failed to get git branch: " + err.Error())
	}
	branch := strings.TrimSpace(string(branchInfo))
	// 如果是主分支，则直接返回
	if strings.EqualFold("master", branch) {
		return
	}
	// 2. 获取当前项目的路径和名称
	currentDir, err := os.Getwd()
	if err != nil {
		panic("Failed to get current directory: " + err.Error())
	}
	splitPath := strings.Split(currentDir, "/")
	projectName := splitPath[len(splitPath)-1]
	projectPath := strings.Join(splitPath[0:len(splitPath)-1], "/")
	// 3. 生成 diff 文件

	dissText := "/Users/knowreason/soft/infer/out/diff/" + projectName + "_diff.txt"
	if !onlyAnalyze {
		command := exec.Command("git", "diff", "--name-only", branch, "master")
		output, err := command.Output()
		if err != nil {
			panic("Failed to get diff: " + strings.Join(command.Args, " ") + " " + err.Error())
		}
		err = os.WriteFile(dissText, output, 0644)
		if err != nil {
			panic("Failed to write diff: " + err.Error())
		}
	}
	// 执行 capture
	captureDir := runInfer(projectPath, projectName, onlyAnalyze, true)
	cmd := exec.Command("infer",
		append(
			append(defaultParams, "-o", captureDir, "analyze", "--changed-files-index", dissText),
		)...,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("完整命令:", strings.Join(cmd.Args, " "))
	if err = cmd.Run(); err != nil {
		panic("Execute infer analyze command failed: " + err.Error())
	}
	fmt.Println("Infer analysis completed. See more execute : infer -o " + captureDir + " explore")
}

func runInfer(projectPath string, projectName string, onlyAnalyze bool, onlyCapture bool) string {
	var fullPath string
	if strings.HasSuffix(projectPath, "/") {
		fullPath = projectPath + projectName
	} else {
		fullPath = projectPath + "/" + projectName
	}
	// 检查是否是个文件夹
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		panic("The project path is not valid, please check it. : " + fullPath)
	}
	if !fileInfo.IsDir() {
		panic("The project path is not a directory.  " + fullPath)
	}
	// 执行 infer 命令
	var inferOutDir = "/Users/knowreason/soft/infer/out/" + projectName
	stat, err := os.Stat(inferOutDir)
	if !onlyAnalyze {
		// 如果他存在，并且是目录
		if err == nil && stat.IsDir() {
			filePath := inferOutDir
			err := exec.Command("rm", "-rf", filePath).Run()
			if err != nil {
				panic("rm file has failed")
			}
		}

		//command := exec.Command("infer",, defaultParams, )
		command := exec.Command("infer",
			append(
				[]string{"--java-version", "17", "capture", "-o", inferOutDir, "--", "mvn", "clean", "compile", "-DskipTests=true"},
			)...,
		)
		command.Dir = fullPath
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		fmt.Println("完整命令:", strings.Join(command.Args, " "))
		if err = command.Run(); err != nil {
			panic("Execute infer command failed: " + err.Error())
		}
	}
	if onlyCapture {
		return inferOutDir
	}

	cmd := exec.Command("infer",
		append(
			append(defaultParams, "-o", inferOutDir, "analyze"),
		)...,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("完整命令:", strings.Join(cmd.Args, " "))
	if err = cmd.Run(); err != nil {
		panic("Execute infer analyze command failed: " + err.Error())
	}
	fmt.Println("Infer analysis completed. See more execute : infer -o " + inferOutDir + " explore")
	return inferOutDir
}

var defaultParams = []string{
	"--annotation-reachability",
	"--bufferoverrun",
	"--config-impact-analysis",
	"--cost",
	"--inefficient-keyset-iterator",
	"--loop-hoisting",
	"--pulse",
	"--racerd",
	//"--resource-leak-lab",
	"--scope-leakage",
	"--starvation",
}
