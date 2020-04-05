package gordon

import (
	"os"
	"path/filepath"
)

func unlinker(target string) func(string, os.FileInfo) error {
	return func(path string, fi os.FileInfo) error {
		if (fi.Mode() & os.ModeSymlink) != os.ModeSymlink {
			return nil
		}
		destination, err := os.Readlink(path)
		if err != nil {
			return err
		}
		if !filepath.HasPrefix(destination, target) {
			return nil
		}
		return os.Remove(path)
	}
}

func Unlink(ev Env, app App) error {
	unlink := unlinker(AppPath(ev, app))
	if err := walkIfDir(ev.Bin(), unlink); err != nil {
		return err
	}
	if err := walkIfDir(ev.Man(), unlink); err != nil {
		return err
	}
	return nil
}
