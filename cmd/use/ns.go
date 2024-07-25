package use

import (
	"github.com/fioncat/ks/pkg/cmdentry"
	"github.com/fioncat/ks/pkg/kubectx"
	"github.com/fioncat/ks/pkg/kubectx/cmdhelper"
	"github.com/fioncat/ks/pkg/metadata"
	"github.com/spf13/cobra"
)

func newNsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ns [NAME]",
		Short: "Switch to a namespace",

		Args: cobra.RangeArgs(0, 1),
	}

	cmdentry.Setup(cmd, RunNs)
	return cmd
}

func RunNs(meta *metadata.Metadata, manager *kubectx.KubeManager, args []string) error {
	ctx, err := manager.GetCurrent()
	if err != nil {
		return err
	}

	ns, err := cmdhelper.SelectNamespace(meta, ctx, args)
	if err != nil {
		return err
	}

	return cmdhelper.PrintUseNamespace(meta, ctx, ns)
}