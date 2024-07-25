package clear

import (
	"github.com/fioncat/ks/pkg/cmdentry"
	"github.com/fioncat/ks/pkg/kubectx"
	"github.com/fioncat/ks/pkg/kubectx/cmdhelper"
	"github.com/fioncat/ks/pkg/metadata"
	"github.com/spf13/cobra"
)

func newNsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ns",
		Short: "Clear the current namespace",

		Args: cobra.NoArgs,
	}

	cmdentry.Setup(cmd, runNs)
	return cmd
}

func runNs(meta *metadata.Metadata, manager *kubectx.KubeManager, args []string) error {
	cmdhelper.PrintClearNamespace()
	return nil
}
