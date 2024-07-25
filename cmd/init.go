package cmd

import (
	"fmt"
	"os"

	"github.com/fioncat/ks/hack"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init [bash|zsh]",
		Short: "Generate init amd completion script",

		DisableFlagsInUseLine: true,

		ValidArgs: []string{"bash", "zsh"},

		Args: cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			root := cmd.Root()
			root.Use = "ks"

			switch args[0] {
			case "bash":
				// same shell script as zsh, but different bash completion
				fmt.Println(hack.Bash)
				return root.GenBashCompletion(os.Stdout)
			case "zsh":
				fmt.Println(hack.Bash)
				return root.GenZshCompletion(os.Stdout)
			}
			return fmt.Errorf("unsupported shell type: %s", args[0])
		},
	}
}
