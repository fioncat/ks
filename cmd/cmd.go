package cmd

import (
	"github.com/fioncat/ks/cmd/use"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "ks",

		SilenceErrors: true,
		SilenceUsage:  true,

		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newEditCmd())
	cmd.AddCommand(use.NewCmd())

	return cmd
}
