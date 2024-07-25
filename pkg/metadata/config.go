package metadata

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	"k8s.io/client-go/util/homedir"
)

const (
	configPathEnvName = "KS_CONFIG"

	recommandedConfigPath  = ".config/ks/config.yaml"
	recommandedMetadataDir = ".local/share/ks"
)

type configPathOptions struct {
	path     string
	explicit bool // If true, the config file MUST exist
}

func newConfigPathOptions(path string) *configPathOptions {
	if path != "" {
		// The user specify config path from command args
		return &configPathOptions{
			path:     path,
			explicit: true,
		}
	}

	path = os.Getenv(configPathEnvName)
	if path != "" {
		// The user specify config path from environment variable
		return &configPathOptions{
			path:     path,
			explicit: true,
		}
	}

	path = filepath.Join(homedir.HomeDir(), recommandedConfigPath)
	return &configPathOptions{
		path:     path,
		explicit: false,
	}
}

type Config struct {
	Groups map[string][]string `yaml:"groups"`
	Alias  map[string]string   `yaml:"alias"`

	MetadataDir string `yaml:"metadataDir"`
}

func newDefaultConfig() *Config {
	return &Config{
		Groups:      make(map[string][]string),
		Alias:       make(map[string]string),
		MetadataDir: filepath.Join(homedir.HomeDir(), recommandedMetadataDir),
	}
}

func loadConfig(configPath string) (*Config, error) {
	pathOptions := newConfigPathOptions(configPath)

	data, err := os.ReadFile(pathOptions.path)
	if err != nil {
		if os.IsNotExist(err) && !pathOptions.explicit {
			return newDefaultConfig(), nil
		}

		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("parse config file %q: %w", pathOptions.path, err)
	}

	// User can add environment variables in metadataDir config, such as "$HOME".
	// Here we need to expand those environment variables.
	cfg.MetadataDir = os.ExpandEnv(cfg.MetadataDir)

	return &cfg, nil
}
