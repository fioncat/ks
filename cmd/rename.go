package cmd

import (
	"fmt"

	"github.com/fioncat/ks/pkg/cmdentry"
	"github.com/fioncat/ks/pkg/kubectx"
	"github.com/fioncat/ks/pkg/kubectx/cmdhelper"
	"github.com/fioncat/ks/pkg/metadata"
	"github.com/spf13/cobra"
)

func newRenameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rename [SRC_NAME] [DST_NAME]",
		Short: "Rename a kubeconfig",

		Args: cobra.ExactArgs(2),

		ValidArgsFunction: cmdhelper.CompleteConfig(false, true),
	}

	cmdentry.Setup(cmd, runRename)
	return cmd
}

func runRename(meta *metadata.Metadata, manager *kubectx.KubeManager, args []string) error {
	srcName, dstName := args[0], args[1]

	ctx, err := manager.Get(srcName)
	if err != nil {
		return err
	}

	_, err = manager.Get(dstName)
	if err == nil {
		return fmt.Errorf("rename destination %q already exists", dstName)
	}

	meta.History.ClearKubeConfig(srcName)
	err = meta.History.Save()
	if err != nil {
		return err
	}

	err = manager.Rename(ctx, dstName)
	if err != nil {
		return err
	}

	if ctx.Current {
		cmdhelper.PrintClearConfig()
	}

	return nil
}
