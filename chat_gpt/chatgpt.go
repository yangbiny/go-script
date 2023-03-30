package main

import (
	"github.com/spf13/cobra"
	"script/chat_gpt/command"
)

var exampleUsage = `
# init for projects
tpctl init -f config.yaml
`

func main() {

	rootCmd := &cobra.Command{
		Use:        "chatGPT",
		Short:      "A Command for using chat_GPT",
		Example:    exampleUsage,
		SuggestFor: []string{"chatGPT"},
	}

	rootCmd.AddCommand(command.ChatCmd())

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
