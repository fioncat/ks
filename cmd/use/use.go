package use

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use",
		Short: "Switch commands",
	}

	cmd.AddCommand(newConfigCmd())
	cmd.AddCommand(newNsCmd())
	cmd.AddCommand(newGroupCmd())

	return cmd
}
