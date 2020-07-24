package gordon

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/saracen/walker"
)

const assetDirName = "asset"

func assetSubPath(ev Env, subs ...string) string {
	return filepath.Join(append([]string{ev.Cache(), assetDirName}, subs...)...)
}

func VersionPath(ev Env, version Version) string {
	return filepath.Join(ev.Cache(), assetDirName, version.owner, version.name, version.tag)
}

func AppPath(ev Env, app App) string {
	return filepath.Join(ev.Cache(), assetDirName, app.owner, app.name)
}

func isDir(dirExpect string) (bool, error) {
	dirStat, err := os.Stat(dirExpect)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if !dirStat.IsDir() {
		return false, nil
	}
	return true, nil
}

func isDirWithChild(dirExpect string) (bool, error) {
	ok, err := isDir(dirExpect)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	children, err := ioutil.ReadDir(dirExpect)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return len(children) > 0, nil
}

func walkIfDir(dirExpect string, walkFn func(pathname string, fi os.FileInfo) error) error {
	ok, err := isDir(dirExpect)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	return walker.Walk(dirExpect, walkFn)
}
