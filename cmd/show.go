package cmd

import (
	"fmt"

	"github.com/fioncat/ks/pkg/cmdentry"
	"github.com/fioncat/ks/pkg/kubectx"
	"github.com/fioncat/ks/pkg/metadata"
	"github.com/spf13/cobra"
)

func newShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show current attached config",

		Args: cobra.NoArgs,
	}

	cmdentry.Setup(cmd, runShow)
	return cmd
}

func runShow(meta *metadata.Metadata, manager *kubectx.KubeManager, args []string) error {
	ctx, err := manager.GetCurrent()
	if err != nil {
		return err
	}

	name := ctx.ConfigName
	if ctx.Alias != "" {
		name = fmt.Sprintf("%s (%s)", name, ctx.Alias)
	}

	if ctx.Namespace == "" {
		fmt.Println(name)
		return nil
	}

	fmt.Printf("%s -> %s\n", name, ctx.Namespace)
	return nil
}
