package gordon

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/blang/semver"
	"github.com/saracen/walker"
)

type Version struct {
	App
	tag string
}

var (
	ErrVersionNotFound = errors.New("version not found")
)

func validateVersionSpec(ev Env, owner, name, tag string) (*Version, error) {
	ver := Version{
		App: App{
			owner: owner,
			name:  name,
		},
		tag: tag,
	}
	path := VersionPath(ev, ver)
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrVersionNotFound
		}
		return nil, err
	}
	if !fi.IsDir() {
		return nil, ErrVersionNotFound
	}
	return &ver, nil
}

func findLatestVersion(ev Env, owner, name string) (*Version, error) {
	var found bool
	var newest semver.Version
	if err := walker.Walk(assetSubPath(ev, owner, name), func(path string, fi os.FileInfo) error {
		if !fi.IsDir() {
			return nil
		}

		tag := filepath.Base(path)
		ver, err := semver.Parse(tag)
		if err != nil {
			// if the tag is not formatted as semver, skip it.
			return filepath.SkipDir
		}
		if ver.GT(newest) {
			newest = ver
			found = true
		}
		return filepath.SkipDir
	}); err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrVersionNotFound
	}
	return &Version{
		App: App{
			owner: owner,
			name:  name,
		},
		tag: newest.String(),
	}, nil
}

func FindVersion(ev Env, spec VersionSpec) (*Version, error) {
	if spec.tag != "" {
		return validateVersionSpec(ev, spec.owner, spec.name, spec.tag)
	}

	return findLatestVersion(ev, spec.owner, spec.name)
}

func ParseVersionPath(ev Env, path string) (*Version, error) {
	rel, err := filepath.Rel(ev.Cache(), path)
	if err != nil {
		return nil, err
	}

	terms := strings.Split(rel, string([]rune{filepath.Separator}))
	if len(terms) < 3 {
		return nil, errors.New("invalid version path")
	}
	if _, err := semver.Parse(terms[2]); err != nil {
		return nil, errors.New("invalid version path")
	}
	ver := Version{
		App: App{
			owner: terms[0],
			name:  terms[1],
		},
		tag: terms[2],
	}
	return &ver, nil
}

type VersionWalker func(Version) error

func WalkVersions(ev Env, walk VersionWalker) error {
	return walker.Walk(assetSubPath(ev), func(path string, fi os.FileInfo) error {
		if !fi.IsDir() {
			return nil
		}
		ver, err := ParseVersionPath(ev, path)
		if err != nil {
			return nil
		}
		if err := walk(*ver); err != nil {
			return err
		}
		return filepath.SkipDir
	})
}

func ListVersions(ev Env) ([]Version, error) {
	var versions []Version
	if err := WalkVersions(ev, func(version Version) error {
		versions = append(versions, version)
		return nil
	}); err != nil {
		return nil, err
	}
	return versions, nil
}

func WalkAppVersions(ev Env, app App, walk VersionWalker) error {
	return walker.Walk(AppPath(ev, app), func(path string, fi os.FileInfo) error {
		if !fi.IsDir() {
			return nil
		}
		tag := filepath.Base(path)
		if _, err := semver.Parse(tag); err != nil {
			// if the tag is not formatted as semver, skip it.
			return filepath.SkipDir
		}
		if err := walk(Version{App: app, tag: tag}); err != nil {
			return err
		}
		return filepath.SkipDir
	})
}

func ListAppVersions(ev Env, app App) ([]Version, error) {
	var versions []Version
	if err := WalkAppVersions(ev, app, func(version Version) error {
		versions = append(versions, version)
		return nil
	}); err != nil {
		return nil, err
	}
	return versions, nil
}

func WalkInstalledVersions(ev Env, walk VersionWalker) error {
	uniq := map[Version]struct{}{}
	if err := walker.Walk(ev.Bin(), walkLinkedVersions(ev, uniq, walk)); err != nil {
		return err
	}
	if err := walker.Walk(ev.Man(), walkLinkedVersions(ev, uniq, walk)); err != nil {
		return err
	}
	return nil
}

func ListInstalledVersions(ev Env) ([]Version, error) {
	var versions []Version
	if err := WalkInstalledVersions(ev, func(version Version) error {
		versions = append(versions, version)
		return nil
	}); err != nil {
		return nil, err
	}
	return versions, nil
}
func walkLinkedVersions(ev Env, uniq map[Version]struct{}, walk func(Version) error) func(path string, fi os.FileInfo) error {
	return func(path string, fi os.FileInfo) error {
		if (fi.Mode() & os.ModeSymlink) != os.ModeSymlink {
			return filepath.SkipDir
		}
		link, err := os.Readlink(path)
		if err != nil {
			return err
		}

		ver, err := ParseVersionPath(ev, link)
		if err != nil {
			return nil // ignore invalid link
		}

		if _, ok := uniq[*ver]; ok {
			return nil
		}
		if err := walk(*ver); err != nil {
			return err
		}
		uniq[*ver] = struct{}{}
		return nil
	}
}
