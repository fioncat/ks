package cmdhelper

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fioncat/ks/pkg/kubectx"
	"github.com/fioncat/ks/pkg/metadata"
	"github.com/fioncat/ks/pkg/utils/fzf"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	kubeClientTimeout = time.Second * 10
)

func SelectContext(meta *metadata.Metadata, manager *kubectx.KubeManager, args []string, require bool) (*kubectx.KubeContext, error) {
	if len(args) == 0 {
		// TODO: We can sort the contexts by use frequency
		ctxs := manager.List()
		if len(ctxs) == 0 {
			return nil, errors.New("no kubeconfig yet, please create one first")
		}
		items := make([]string, 0, len(ctxs))
		for _, ctx := range ctxs {
			items = append(items, ctx.ConfigName)
		}

		idx, err := fzf.Search(items)
		if err != nil {
			return nil, err
		}

		return ctxs[idx], nil
	}

	name := args[0]
	if name == "-" {
		currentCtx, _ := manager.GetCurrent()
		var currentConfigName string
		if currentCtx != nil {
			currentConfigName = currentCtx.ConfigName
		}

		lastName := meta.History.GetLastKubeConfig(currentConfigName)
		if lastName == nil {
			return nil, errors.New("no last switched kubeconfig")
		}

		name = *lastName
	}

	ctx, err := manager.Get(name)
	if err != nil {
		if !require {
			return nil, nil
		}
		return nil, err
	}

	return ctx, nil
}

func ListNamespaces(ctx *kubectx.KubeContext) ([]string, error) {
	config, err := clientcmd.BuildConfigFromFlags("", ctx.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("load kubeconfig: %w", err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("create kubernetes client: %w", err)
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), kubeClientTimeout)
	defer cancel()
	nsList, err := client.CoreV1().Namespaces().List(timeoutCtx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list namespaces: %w", err)
	}

	namespaces := make([]string, 0, len(nsList.Items))
	for _, nsItem := range nsList.Items {
		namespaces = append(namespaces, nsItem.Name)
	}

	return namespaces, nil
}

func SelectNamespace(meta *metadata.Metadata, ctx *kubectx.KubeContext, args []string) (string, error) {
	if len(args) > 0 {
		name := args[0]
		if name == "-" {
			lastNamespace := meta.History.GetLastNamespace(ctx.ConfigName, ctx.Namespace)
			if lastNamespace == nil {
				return "", errors.New("no last switched namespace")
			}
			name = *lastNamespace
		}
		return name, nil
	}

	namespaces, err := ListNamespaces(ctx)
	if err != nil {
		return "", err
	}
	if len(namespaces) == 0 {
		return "", errors.New("no namespace in the cluster to use")
	}

	idx, err := fzf.Search(namespaces)
	if err != nil {
		return "", err
	}

	return namespaces[idx], nil
}
