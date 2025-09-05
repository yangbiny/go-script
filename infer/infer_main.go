package main

import (
	"script/infer/command"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:        "pinfer",
		Short:      "Command for facebook infer",
		Example:    "infer run [path to project]",
		SuggestFor: []string{""},
	}
	rootCmd.AddCommand(infer.RunInfer())
	rootCmd.Execute()
}
