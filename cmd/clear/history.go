package clear

import (
	"github.com/fioncat/ks/pkg/cmdentry"
	"github.com/fioncat/ks/pkg/kubectx"
	"github.com/fioncat/ks/pkg/metadata"
	"github.com/spf13/cobra"
)

func newHistoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history",
		Short: "Clear all history",

		Args: cobra.NoArgs,
	}

	cmdentry.Setup(cmd, runHistory)
	return cmd
}

func runHistory(meta *metadata.Metadata, manager *kubectx.KubeManager, args []string) error {
	meta.History.ClearAll()
	return meta.History.Save()
}
