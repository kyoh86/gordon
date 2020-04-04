package gordon

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/saracen/walker"
)

const executable = 0400

func unlinker(target string) func(string, os.FileInfo) error {
	return func(path string, fi os.FileInfo) error {
		if (fi.Mode() & os.ModeSymlink) != os.ModeSymlink {
			return nil
		}
		destination, err := os.Readlink(path)
		if err != nil {
			return err
		}
		if filepath.HasPrefix(destination, target) {
			return nil
		}
		return os.Remove(path)
	}
}

func Unlink(ev Env, app *App) error {
	unlink := unlinker(AppPath(ev, app))
	if err := walker.Walk(ev.Bin(), unlink); err != nil {
		return err
	}
	if err := walker.Walk(ev.Man(), unlink); err != nil {
		return err
	}
	return nil
}

func Link(ev Env, release *Release) error {
	bins := map[string]string{}
	mans := map[string]string{}
	// unlink old links
	if err := Unlink(ev, &release.App); err != nil {
		return err
	}

	// link all executables and mans
	if err := walker.Walk(ReleasePath(ev, release), func(path string, fi os.FileInfo) error {
		switch {
		case (fi.Mode() & executable) == executable:
			// executable file
			bins[path] = fi.Name()
		case strings.HasSuffix(path, ".1"):
			// man file
			mans[path] = fi.Name()
		}
		return nil
	}); err != nil {
		return err
	}

	for binPath, binName := range bins {
		if err := os.Link(binPath, filepath.Join(ev.Bin(), binName)); err != nil {
			return err
		}

		manPath := binPath + ".1"
		manName, ok := mans[manPath]
		if !ok {
			continue
		}

		if err := os.Link(manPath, filepath.Join(ev.Man(), manName)); err != nil {
			return err
		}
	}
	return nil
}
