package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/fioncat/ks/cmd"
	"github.com/fioncat/ks/pkg/utils/fzf"
)

var (
	Version   string = "N/A"
	Commit    string = "N/A"
	BuildDate string = "N/A"
)

func main() {
	cmd := cmd.NewCmd()

	cmd.Version = Version
	err := cmd.Execute()
	if err != nil {
		if errors.Is(err, fzf.ErrCanceled) {
			os.Exit(fzf.ExitCodeCanceled)
		}

		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
