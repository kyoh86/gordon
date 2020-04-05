package command

import (
	"fmt"
	"io"
	"os"

	"github.com/kyoh86/gordon/internal/gordon"
)

func Dump(ev Env, bundleFile string) error {
	var writer io.Writer
	if bundleFile == "-" {
		writer = os.Stdout
	} else {
		f, err := os.OpenFile(bundleFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		writer = f
		defer f.Close()
	}
	if err := gordon.WalkInstalledVersions(ev, func(ver gordon.Version) error {
		fmt.Fprintln(writer, ver)
		return nil
	}); err != nil {
		return err
	}
	return nil
}
