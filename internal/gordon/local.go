package gordon

import (
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

func walkIfDir(dirExpect string, walkFn func(pathname string, fi os.FileInfo) error) error {
	dirStat, err := os.Stat(dirExpect)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if !dirStat.IsDir() {
		return nil
	}
	return walker.Walk(dirExpect, walkFn)
}
