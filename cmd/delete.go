package cmd

import (
	"fmt"

	"github.com/fioncat/ks/pkg/cmdentry"
	"github.com/fioncat/ks/pkg/kubectx"
	"github.com/fioncat/ks/pkg/kubectx/cmdhelper"
	"github.com/fioncat/ks/pkg/metadata"
	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [NAME]",
		Short: "Delete a kubeconfig",

		Args: cobra.RangeArgs(0, 1),

		ValidArgsFunction: cmdhelper.CompleteConfig(false),
	}

	cmdentry.Setup(cmd, runDelete)
	return cmd
}

func runDelete(meta *metadata.Metadata, manager *kubectx.KubeManager, args []string) error {
	ctx, err := cmdhelper.SelectContext(meta, manager, args, true, false)
	if err != nil {
		return err
	}

	meta.History.ClearKubeConfig(ctx.ConfigName)
	err = meta.History.Save()
	if err != nil {
		return fmt.Errorf("clear history for config: %w", err)
	}

	err = manager.Delete(ctx.ConfigName)
	if err != nil {
		return err
	}

	if ctx.Current {
		cmdhelper.PrintClearConfig()
	}

	return nil
}
