package cmdhelper

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/fioncat/ks/pkg/kubectx"
	"github.com/fioncat/ks/pkg/metadata"
	"github.com/spf13/cobra"
)

func CompleteConfig(skipCurrent, rename bool) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 1 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		if len(args) == 1 && !rename {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		manager := buildCompleteKubeManager(args, toComplete)
		if manager == nil {
			return nil, cobra.ShellCompDirectiveError
		}

		items, err := completeConfig(manager, skipCurrent)
		if err != nil {
			handleCompleteError(args, toComplete, err)
			return nil, cobra.ShellCompDirectiveError
		}

		if rename && len(args) == 1 {
			return items, cobra.ShellCompDirectiveNoSpace
		}

		return items, cobra.ShellCompDirectiveNoFileComp
	}
}

func completeConfig(manager *kubectx.KubeManager, skipCurrent bool) ([]string, error) {
	ctxs := manager.List()
	items := make([]string, 0, len(ctxs))
	for _, ctx := range ctxs {
		if ctx.Current && skipCurrent {
			continue
		}
		items = append(items, ctx.ConfigName)
	}

	return items, nil
}

func CompleteNamespace(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	manager := buildCompleteKubeManager(args, toComplete)
	if manager == nil {
		return nil, cobra.ShellCompDirectiveError
	}

	items, err := completeNamespace(manager)
	if err != nil {
		handleCompleteError(args, toComplete, err)
		return nil, cobra.ShellCompDirectiveError
	}

	return items, cobra.ShellCompDirectiveNoFileComp
}

func CompleteGroupNamespace(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	meta, err := metadata.Load("")
	if err != nil {
		handleCompleteError(args, toComplete, err)
		return nil, cobra.ShellCompDirectiveError
	}

	if len(meta.Config.Groups) == 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if len(args) == 0 {
		groups := make([]string, 0, len(meta.Config.Groups))
		for group := range meta.Config.Groups {
			groups = append(groups, group)
		}
		sort.Strings(groups)
		return groups, cobra.ShellCompDirectiveNoFileComp
	}

	if len(args) != 1 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	group := args[0]
	items := meta.Config.Groups[group]
	return items, cobra.ShellCompDirectiveNoFileComp
}

func completeNamespace(manager *kubectx.KubeManager) ([]string, error) {
	ctx, err := manager.GetCurrent()
	if err != nil {
		return nil, err
	}

	items, err := ListNamespaces(ctx, ctx.Namespace)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func buildCompleteKubeManager(args []string, toComplete string) *kubectx.KubeManager {
	// TODO: Parse args to check if the `--config` option is set.
	// In that cases, we should read the specified config file.
	meta, err := metadata.Load("")
	if err != nil {
		handleCompleteError(args, toComplete, err)
		return nil
	}

	manager, err := kubectx.BuildKubeManager(meta)
	if err != nil {
		handleCompleteError(args, toComplete, err)
		return nil
	}

	return manager
}

func handleCompleteError(args []string, toComplete string, err error) {
	msg := fmt.Sprintf("Complete error: args=%v, toComplete=%q, error=%v\n", args, toComplete, err)

	logsFilePath := filepath.Join(os.TempDir(), "ks_complete_error.log")
	file, err := os.OpenFile(logsFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open complete logs file: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write complete logs file: %v\n", err)
		return
	}
}
