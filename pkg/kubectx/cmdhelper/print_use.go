package cmdhelper

import (
	"fmt"
	"strings"

	"github.com/fioncat/ks/pkg/kubectx"
	"github.com/fioncat/ks/pkg/metadata"
)

const (
	ConfigHackPrefix    = "__config_"
	NamespaceHackPrefix = "__namespace_"
)

func PrintUseConfig(meta *metadata.Metadata, ctx *kubectx.KubeContext) error {
	fields := []string{
		ctx.ConfigName,
		ctx.ConfigPath,
	}

	meta.History.Add(ctx.ConfigName, "")
	err := meta.History.Save()
	if err != nil {
		return err
	}

	fmt.Printf("%s%s\n", ConfigHackPrefix, strings.Join(fields, ","))
	return nil
}

func PrintUseNamespace(meta *metadata.Metadata, ctx *kubectx.KubeContext, ns string) error {
	ctx.Namespace = ns
	meta.History.Add(ctx.ConfigName, ns)
	err := meta.History.Save()
	if err != nil {
		return err
	}

	fmt.Printf("%s%s\n", NamespaceHackPrefix, ns)
	return nil
}

func PrintClearConfig() {
	fmt.Printf("%s,\n", ConfigHackPrefix)
}

func PrintClearNamespace() {
	fmt.Println(NamespaceHackPrefix)
}
