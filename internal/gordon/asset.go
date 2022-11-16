package gordon

import (
	"strings"

	"github.com/google/go-github/v29/github"
	"github.com/kyoh86/gordon/internal/archive"
)

func findAsset(ev Env, assets []github.ReleaseAsset) (asset github.ReleaseAsset, _ error) {
	matchers := []func(s string) bool{
		func(s string) bool {
			return matchOpener(s) != nil
		},
	}

	arch, ok := archMatches[ev.Architecture()]
	if !ok {
		arch = func(s string) bool {
			return strings.Contains(s, ev.Architecture())
		}
	}
	matchers = append(matchers, arch)

	os, ok := osMatches[ev.OS()]
	if !ok {
		os = func(s string) bool {
			return strings.Contains(s, ev.OS())
		}
	}
	matchers = append(matchers, os)

	for _, matcher := range matchers {
		if len(assets) == 1 {
			return assets[0], nil
		}
		assets = matchAssets(assets, matcher)
	}
	if len(assets) == 0 {
		return asset, ErrAssetNotFound
	}
	return assets[0], nil
}

func matchAssets(assets []github.ReleaseAsset, matcher func(s string) bool) []github.ReleaseAsset {
	var res []github.ReleaseAsset
	for _, asset := range assets {
		if matcher(asset.GetName()) {
			res = append(res, asset)
		}
	}
	if len(res) == 0 {
		return assets
	}
	return res
}

var openers = map[string]archive.Opener{
	".tar.xz":  archive.OpenTarXz,
	".tgz":     archive.OpenTarGzip,
	".tar.gz":  archive.OpenTarGzip,
	".tar.bz2": archive.OpenTarBzip2,
	".tar":     archive.OpenTar,
	".zip":     archive.OpenZip,
}

func matchOpener(name string) archive.Opener {
	ln := strings.ToLower(name)
	for suffix, opener := range openers {
		if strings.HasSuffix(ln, suffix) {
			return opener
		}
	}
	return nil
}

func containsOneOf(s string, matchList ...string) bool {
	for _, m := range matchList {
		if strings.Contains(s, m) {
			return true
		}
	}
	return false
}

func match386(s string) bool {
	if strings.Contains(strings.ToLower(s), "x86_64") {
		return false
	}
	return containsOneOf(strings.ToLower(s), "386", "686", "linux32", "x86")
}

func matchAMD64(s string) bool {
	return containsOneOf(strings.ToLower(s), "x86_64", "amd64", "intel", "linux64")
}

func matchARM64(s string) bool {
	return containsOneOf(strings.ToLower(s), "aarch64", "arm64")
}

var archMatches = map[string]func(string) bool{
	"386":   match386,
	"amd64": matchAMD64,
	"arm64": matchARM64,
}

func matchLinux(s string) bool {
	return strings.Contains(strings.ToLower(s), "linux")
}

func matchDarwin(s string) bool {
	return containsOneOf(strings.ToLower(s), "darwin", "mac", "osx", "os-x")
}

func matchWindows(s string) bool {
	return containsOneOf(strings.ToLower(s), "windows|-win|_win|win64|win32")
}

var osMatches = map[string]func(s string) bool{
	"linux":   matchLinux,
	"darwin":  matchDarwin,
	"windows": matchWindows,
}
