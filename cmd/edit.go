package cmd

import (
	"github.com/fioncat/ks/pkg/cmdentry"
	"github.com/fioncat/ks/pkg/kubectx"
	"github.com/fioncat/ks/pkg/kubectx/cmdhelper"
	"github.com/fioncat/ks/pkg/metadata"
	"github.com/fioncat/ks/pkg/utils"
	"github.com/spf13/cobra"
)

func newEditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit [NAME]",
		Short: "Edit a kubeconfig file",

		Args: cobra.RangeArgs(0, 1),

		ValidArgsFunction: cmdhelper.CompleteConfig(false, false),
	}

	cmdentry.Setup(cmd, runEdit)
	return cmd
}

func runEdit(meta *metadata.Metadata, manager *kubectx.KubeManager, args []string) error {
	ctx, err := cmdhelper.SelectContext(meta, manager, args, false, false)
	if err != nil {
		return err
	}

	var sourcePath *string
	var name string
	if ctx != nil {
		sourcePath = &ctx.ConfigPath
		name = ctx.ConfigName
	} else {
		name = args[0]
	}

	content, err := utils.Edit(meta.TempDir, sourcePath)
	if err != nil {
		return err
	}

	_, err = manager.Set(name, content)
	if err != nil {
		return err
	}

	return nil
}
