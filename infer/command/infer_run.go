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
			runInfer(projectPath, projectName, onlyAnalyze)
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

func runInfer(projectPath string, projectName string, onlyAnalyze bool) {
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
