package main

import (
	"github.com/spf13/cobra"
	utils "script/utils/command"
)

var exampleUsage = `
	utils file csv2excel
`

func main() {

	rootCmd := &cobra.Command{
		Use:        "utils",
		Short:      "常用工具，各种转换工具的实现",
		Example:    exampleUsage,
		SuggestFor: []string{""},
	}

	rootCmd.AddCommand(utils.FileCommand())
	rootCmd.AddCommand(utils.ArmeriaDebugCommand())
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}

}
