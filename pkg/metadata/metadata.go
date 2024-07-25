package metadata

import (
	"fmt"
	"path/filepath"

	"github.com/fioncat/ks/pkg/utils"
)

const (
	historyFileName   = "history.yaml"
	kubeConfigDirName = "configs"
)

type Metadata struct {
	Config *Config

	History *History

	KubeConfigDir string

	TempDir string
}

func Load(configPath string) (*Metadata, error) {
	config, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}

	historyPath := filepath.Join(config.MetadataDir, historyFileName)
	history, err := loadHistory(historyPath)
	if err != nil {
		return nil, err
	}

	kubeConfigDir := filepath.Join(config.MetadataDir, kubeConfigDirName)
	err = utils.EnsureDir(kubeConfigDir)
	if err != nil {
		return nil, fmt.Errorf("ensure kube configs dir: %w", err)
	}

	return &Metadata{
		Config:        config,
		History:       history,
		KubeConfigDir: kubeConfigDir,
		TempDir:       filepath.Join(config.MetadataDir, "tmp"),
	}, nil
}
