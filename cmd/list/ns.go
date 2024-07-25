package list

import (
	"github.com/fioncat/ks/pkg/kubectx"
	"github.com/fioncat/ks/pkg/kubectx/cmdhelper"
	"github.com/fioncat/ks/pkg/metadata"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type nsItem struct {
	Current   bool   `json:"current"`
	Namespace string `json:"namespace"`
}

func newNsCmd() *cobra.Command {
	return buildCmd("ns", table.Row{
		" ", "NAMESPACE",
	}, func(meta *metadata.Metadata, manager *kubectx.KubeManager) ([]*nsItem, error) {
		ctx, err := manager.GetCurrent()
		if err != nil {
			return nil, err
		}
		namespaces, err := cmdhelper.ListNamespaces(ctx, "")
		if err != nil {
			return nil, err
		}

		items := make([]*nsItem, 0, len(namespaces))
		for _, ns := range namespaces {
			current := ns == ctx.Namespace
			items = append(items, &nsItem{
				Current:   current,
				Namespace: ns,
			})
		}

		return items, nil

	}, func(item *nsItem) table.Row {
		current := ""
		if item.Current {
			current = "*"
		}

		return table.Row{
			current,
			item.Namespace,
		}
	})
}
