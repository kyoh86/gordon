package gordon

import (
	"path/filepath"
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
