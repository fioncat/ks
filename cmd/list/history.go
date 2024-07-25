package list

import (
	"github.com/fioncat/ks/pkg/kubectx"
	"github.com/fioncat/ks/pkg/metadata"
	"github.com/fioncat/ks/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func newHistoryCmd() *cobra.Command {
	return buildCmd("history", table.Row{
		"NAME", "NAMESPACE", "TIME",
	}, func(meta *metadata.Metadata, manager *kubectx.KubeManager) ([]*metadata.HistoryRecord, error) {
		return meta.History.Records, nil
	}, func(record *metadata.HistoryRecord) table.Row {
		return table.Row{
			record.Name,
			record.Namespace,
			utils.FormatTime(record.Timestamp),
		}
	})
}
