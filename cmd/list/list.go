package list

import (
	"encoding/json"
	"fmt"

	"github.com/fioncat/ks/pkg/cmdentry"
	"github.com/fioncat/ks/pkg/kubectx"
	"github.com/fioncat/ks/pkg/metadata"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List commands",
	}

	cmd.AddCommand(newConfigCmd())
	cmd.AddCommand(newHistoryCmd())
	cmd.AddCommand(newNsCmd())

	return cmd
}

type fetchItemsFunc[T any] func(meta *metadata.Metadata, manager *kubectx.KubeManager) ([]*T, error)

type getRowFunc[T any] func(item *T) table.Row

func buildCmd[T any](name string, title table.Row, fetchItems fetchItemsFunc[T], getRow getRowFunc[T]) *cobra.Command {
	var jsonFormat bool
	cmd := &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("List %s", name),

		Args: cobra.NoArgs,
	}

	cmd.Flags().BoolVarP(&jsonFormat, "json", "j", false, "output as json")

	cmdentry.Setup(cmd, func(meta *metadata.Metadata, manager *kubectx.KubeManager, args []string) error {
		items, err := fetchItems(meta, manager)
		if err != nil {
			return err
		}
		return run(name, items, title, getRow, jsonFormat)
	})

	return cmd
}

func run[T any](name string, items []*T, title table.Row, getRow func(item *T) table.Row, jsonFormat bool) error {
	if len(items) == 0 {
		if jsonFormat {
			fmt.Println("[]")
			return nil
		}
		return fmt.Errorf("no %s to list", name)
	}

	if jsonFormat {
		data, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	t := table.NewWriter()
	t.AppendHeader(title)

	for _, item := range items {
		row := getRow(item)
		t.AppendRow(row)
	}

	fmt.Printf("TOTAL: %d\n", len(items))
	fmt.Println(t.Render())
	return nil
}
