package use

import (
	"github.com/fioncat/ks/pkg/cmdentry"
	"github.com/fioncat/ks/pkg/kubectx"
	"github.com/fioncat/ks/pkg/kubectx/cmdhelper"
	"github.com/fioncat/ks/pkg/metadata"
	"github.com/spf13/cobra"
)

func newGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group GROUP [NAME]",
		Short: "Switch to a group namespace",

		Args: cobra.RangeArgs(1, 2),

		ValidArgsFunction: cmdhelper.CompleteGroupNamespace,
	}

	cmdentry.Setup(cmd, runGroup)
	return cmd
}

func runGroup(meta *metadata.Metadata, manager *kubectx.KubeManager, args []string) error {
	ns, err := cmdhelper.SelectGroupNamespace(meta, args)
	if err != nil {
		return err
	}

	ctx, err := manager.GetCurrent()
	if err != nil {
		return err
	}

	return cmdhelper.PrintUseNamespace(meta, ctx, ns)
}
