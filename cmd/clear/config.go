package clear

import (
	"github.com/fioncat/ks/pkg/cmdentry"
	"github.com/fioncat/ks/pkg/kubectx"
	"github.com/fioncat/ks/pkg/kubectx/cmdhelper"
	"github.com/fioncat/ks/pkg/metadata"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Clear the current kubeconfig",

		Args: cobra.NoArgs,
	}

	cmdentry.Setup(cmd, runConfig)
	return cmd
}

func runConfig(meta *metadata.Metadata, manager *kubectx.KubeManager, args []string) error {
	cmdhelper.PrintClearConfig()
	return nil
}
