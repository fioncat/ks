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
	fmt.Printf("%s%s\n", ConfigHackPrefix, strings.Join(fields, ","))

	meta.History.Add(ctx.ConfigName, "")
	return meta.History.Save()
}

func PrintUseNamespace(meta *metadata.Metadata, ctx *kubectx.KubeContext, ns string) error {
	fmt.Printf("%s%s\n", NamespaceHackPrefix, ns)

	ctx.Namespace = ns
	meta.History.Add(ctx.ConfigName, ns)
	return meta.History.Save()
}

func PrintClearConfig() {
	fmt.Printf("%s,\n", ConfigHackPrefix)
}

func PrintClearNamespace() {
	fmt.Println(NamespaceHackPrefix)
}
