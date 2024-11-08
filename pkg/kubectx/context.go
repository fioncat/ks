package kubectx

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fioncat/ks/pkg/metadata"
	"github.com/fioncat/ks/pkg/utils"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type KubeContext struct {
	ConfigName string `json:"configName"`
	ConfigPath string `json:"configPath"`

	Namespace string `json:"namespace"`

	Alias string `json:"alias"`

	Current bool `json:"current"`

	config       *clientcmdapi.Config
	configAccess clientcmd.ConfigAccess
}

func write(meta *metadata.Metadata, name string, data []byte) (*KubeContext, error) {
	path := filepath.Join(meta.KubeConfigDir, name)
	err := utils.WriteFile(path, data)
	if err != nil {
		return nil, err
	}

	ctx, err := read(meta, path)
	if err != nil {
		removeErr := os.Remove(path)
		if removeErr != nil {
			return nil, fmt.Errorf("remove bad kubeconfig after checking: %w", removeErr)
		}

		return nil, fmt.Errorf("invalid kubeconfig content to write: %w", err)
	}

	return ctx, nil
}

func scan(meta *metadata.Metadata) ([]*KubeContext, error) {
	var ctxs []*KubeContext
	err := filepath.Walk(meta.KubeConfigDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ctx, err := read(meta, path)
		if err != nil {
			return fmt.Errorf("read kubeconfig file %q: %w", path, err)
		}

		ctxs = append(ctxs, ctx)
		return nil
	})
	return ctxs, err
}

func read(meta *metadata.Metadata, path string) (*KubeContext, error) {
	name, err := filepath.Rel(meta.KubeConfigDir, path)
	if err != nil {
		// This error should not happen in normal case. So we add a prefix to mark it.
		// If this happens, it means there is a bug in the code.
		return nil, fmt.Errorf("[Internal] bad kubeconfig path %q, not in expected position", path)
	}

	err = utils.EnsureFile(path)
	if err != nil {
		return nil, err
	}

	configAccess := clientcmd.NewDefaultPathOptions()
	// Use ExplicitPath to avoid loading the default kubeconfig if file not exists
	configAccess.LoadingRules.ExplicitPath = path

	config, err := configAccess.GetStartingConfig()
	if err != nil {
		return nil, err
	}

	var namespace string
	// Get namespace from kubeconfig context setting
	ctx := config.Contexts[config.CurrentContext]
	if ctx != nil {
		namespace = ctx.Namespace
	}

	return &KubeContext{
		ConfigName:   name,
		ConfigPath:   path,
		Namespace:    namespace,
		Alias:        "",
		config:       config,
		configAccess: configAccess,
	}, nil
}

func remove(meta *metadata.Metadata, name string) error {
	path := filepath.Join(meta.KubeConfigDir, name)
	return utils.RemoveFileRecursively(path, meta.KubeConfigDir)
}

func getAlias(name string, target *KubeContext) *KubeContext {
	return &KubeContext{
		ConfigName:   name,
		ConfigPath:   target.ConfigPath,
		Namespace:    target.Namespace,
		Alias:        target.ConfigName,
		Current:      target.Current,
		config:       target.config,
		configAccess: target.configAccess,
	}
}
