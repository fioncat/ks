package clear

import "github.com/spf13/cobra"

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear commands",
	}

	cmd.AddCommand(newConfigCmd())
	cmd.AddCommand(newNsCmd())
	cmd.AddCommand(newHistoryCmd())

	return cmd
}
