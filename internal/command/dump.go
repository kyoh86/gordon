package command

import (
	"fmt"
	"io"
	"os"
	"sort"

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
	versions, err := gordon.ListInstalledVersions(ev)
	if err != nil {
		return err
	}
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].String() < versions[j].String()
	})
	for _, ver := range versions {
		fmt.Fprintln(writer, ver)
	}
	return nil
}
