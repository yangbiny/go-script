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
		Short:      "A Command for Common utils",
		Example:    exampleUsage,
		SuggestFor: []string{"file"},
	}

	rootCmd.AddCommand(utils.FileCommand())
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}

}
