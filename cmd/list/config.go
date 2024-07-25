package list

import (
	"github.com/fioncat/ks/pkg/kubectx"
	"github.com/fioncat/ks/pkg/metadata"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	return buildCmd("config", table.Row{
		" ", "NAME", "NAMESPACE", "ALIAS",
	}, func(meta *metadata.Metadata, manager *kubectx.KubeManager) ([]*kubectx.KubeContext, error) {
		return manager.List(), nil
	}, func(ctx *kubectx.KubeContext) table.Row {
		current := ""
		if ctx.Current {
			current = "*"
		}

		return table.Row{
			current,
			ctx.ConfigName,
			ctx.Namespace,
			ctx.Alias,
		}
	})
}
