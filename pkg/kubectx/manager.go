package kubectx

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/fioncat/ks/pkg/metadata"
)

const (
	CurrentKubeConfigNameEnv = "KS_CURRENT_KUBECONFIG_NAME"
	CurrentNamespaceEnv      = "KS_CURRENT_NAMESPACE"
)

type KubeManager struct {
	meta *metadata.Metadata

	ctxs  []*KubeContext
	index map[string]*KubeContext
}

func BuildKubeManager(meta *metadata.Metadata) (*KubeManager, error) {
	ctxs, err := scan(meta)
	if err != nil {
		return nil, fmt.Errorf("scan kubeconfig files: %w", err)
	}

	index := make(map[string]*KubeContext, len(ctxs))
	for _, ctx := range ctxs {
		index[ctx.ConfigName] = ctx
	}

	if len(meta.Config.Alias) > 0 {
		for aliasConfigName, targetConfigName := range meta.Config.Alias {
			if _, ok := index[aliasConfigName]; ok {
				return nil, fmt.Errorf("alias %q conflicts with a kubeconfig file", aliasConfigName)
			}
			targetConfig := index[targetConfigName]
			if targetConfig == nil {
				return nil, fmt.Errorf("alias %q points to a non-exist kubeconfig %q", aliasConfigName, targetConfigName)
			}

			aliasContext := getAlias(aliasConfigName, targetConfig)
			ctxs = append(ctxs, aliasContext)
			index[aliasConfigName] = aliasContext
		}
	}

	sort.Slice(ctxs, func(i, j int) bool {
		return ctxs[i].ConfigName < ctxs[j].ConfigName
	})

	currentName := os.Getenv(CurrentKubeConfigNameEnv)
	currentNamespace := os.Getenv(CurrentNamespaceEnv)
	for _, ctx := range ctxs {
		ctx.Current = currentName == ctx.ConfigName
		if ctx.Current && currentNamespace != "" {
			ctx.Namespace = currentNamespace
		}
	}

	return &KubeManager{
		meta:  meta,
		ctxs:  ctxs,
		index: index,
	}, nil
}

func (m *KubeManager) Set(name string, data []byte) (*KubeContext, error) {
	ctx, err := write(m.meta, name, data)
	if err != nil {
		return nil, err
	}

	var setIndex int = -1
	for idx, ctx := range m.ctxs {
		if ctx.ConfigName == name {
			setIndex = idx
			break
		}
	}

	if setIndex >= 0 {
		m.ctxs[setIndex] = ctx
	}
	m.index[name] = ctx
	return ctx, nil
}

func (m *KubeManager) List() []*KubeContext {
	return m.ctxs
}

func (m *KubeManager) Get(name string) (*KubeContext, error) {
	ctx, ok := m.index[name]
	if !ok {
		return nil, fmt.Errorf("kubeconfig %q not found", name)
	}
	return ctx, nil
}

func (m *KubeManager) Delete(name string) error {
	ctx, ok := m.index[name]
	if !ok {
		return fmt.Errorf("kubeconfig %q not found", name)
	}

	err := remove(m.meta, ctx.ConfigName)
	if err != nil {
		return fmt.Errorf("remove kubeconfig %q: %w", name, err)
	}

	delete(m.index, name)
	for idx, c := range m.ctxs {
		if c.ConfigName == name {
			m.ctxs = append(m.ctxs[:idx], m.ctxs[idx+1:]...)
			break
		}
	}

	return nil
}

func (m *KubeManager) GetCurrent() (*KubeContext, error) {
	for _, ctx := range m.ctxs {
		if ctx.Current {
			return ctx, nil
		}
	}
	return nil, errors.New("no current kubeconfig attached, please attach one first (`use` command)")
}

func (m *KubeManager) Rename(ctx *KubeContext, name string) error {
	newPath := filepath.Join(m.meta.KubeConfigDir, name)
	err := os.Rename(ctx.ConfigPath, newPath)
	if err != nil {
		return fmt.Errorf("rename kubeconfig file: %w", err)
	}

	ctx.ConfigName = name
	ctx.ConfigPath = newPath

	return nil
}
