package cmdentry

import (
	"github.com/fioncat/ks/pkg/kubectx"
	"github.com/fioncat/ks/pkg/metadata"
	"github.com/spf13/cobra"
)

type Runner func(meta *metadata.Metadata, manager *kubectx.KubeManager, args []string) error

func Setup(cmd *cobra.Command, runner Runner) {
	var configPath string
	cmd.Flags().StringVarP(&configPath, "config", "", "", "the path of config file, default is ~/.config/ks/config.yaml")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		meta, err := metadata.Load(configPath)
		if err != nil {
			return err
		}

		manager, err := kubectx.BuildKubeManager(meta)
		if err != nil {
			return err
		}

		return runner(meta, manager, args)
	}
}
