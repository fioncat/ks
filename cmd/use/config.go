package use

import (
	"github.com/fioncat/ks/pkg/cmdentry"
	"github.com/fioncat/ks/pkg/kubectx"
	"github.com/fioncat/ks/pkg/kubectx/cmdhelper"
	"github.com/fioncat/ks/pkg/metadata"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config [NAME]",
		Short: "Switch to a kubeconfig",

		Args: cobra.RangeArgs(0, 1),
	}

	cmdentry.Setup(cmd, runConfig)
	return cmd
}

func runConfig(meta *metadata.Metadata, manager *kubectx.KubeManager, args []string) error {
	ctx, err := cmdhelper.SelectContext(meta, manager, args, true)
	if err != nil {
		return err
	}

	return cmdhelper.PrintUseConfig(meta, ctx)
}
