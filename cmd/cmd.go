package cmd

import (
	"github.com/fioncat/ks/cmd/clear"
	"github.com/fioncat/ks/cmd/list"
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
	cmd.AddCommand(newDeleteCmd())
	cmd.AddCommand(newShowCmd())
	cmd.AddCommand(use.NewCmd())
	cmd.AddCommand(clear.NewCmd())
	cmd.AddCommand(list.NewCmd())

	return cmd
}
