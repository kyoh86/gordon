package gordon

import (
	"os"
	"path/filepath"
	"strings"
)

func Link(ev Env, version Version) error {
	// unlink old links
	if err := Unlink(ev, version.App); err != nil {
		return err
	}

	verDir := VersionPath(ev, version)
	bins := map[string]string{}
	mans := map[string]string{}
	// link all executables and mans
	if err := walkIfDir(verDir, func(path string, fi os.FileInfo) error {
		rel, err := filepath.Rel(verDir, path)
		if err != nil {
			return err
		}
		if rel != fi.Name() {
			return nil
		}
		switch {
		case isExecutable(fi):
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

	if len(bins) == 0 {
		return nil
	}

	if err := os.MkdirAll(ev.Bin(), 0777); err != nil {
		return err
	}
	if err := os.MkdirAll(ev.Man(), 0777); err != nil {
		return err
	}
	for binPath, binName := range bins {
		if err := os.Symlink(binPath, filepath.Join(ev.Bin(), binName)); err != nil {
			return err
		}

		manPath := binPath + ".1"
		manName, ok := mans[manPath]
		if !ok {
			continue
		}

		if err := os.Symlink(manPath, filepath.Join(ev.Man(), manName)); err != nil {
			return err
		}
	}

	return nil
}
