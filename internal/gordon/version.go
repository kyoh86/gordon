package gordon

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/blang/semver"
)

type Version struct {
	App
	tag    string
	semver semver.Version
}

func (v Version) String() string {
	return fmt.Sprintf("%s/%s@%s", v.owner, v.name, v.tag)
}

func (v Version) Tag() string {
	return v.tag
}

func (v Version) Semver() semver.Version {
	return v.semver
}

func (v Version) Spec() VersionSpec {
	return VersionSpec{
		AppSpec:AppSpec{
			owner: v.owner,
			name:  v.name,
		},
		raw    : v.String(),
		tag    : v.tag,
		semver : v.semver,
	}
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
	sv, err := ValidateTag(tag)
	if err != nil {
		return nil, err
	}
	ver.semver = sv
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
	var newest Version
	if err := walkIfDir(assetSubPath(ev, owner, name), func(path string, fi os.FileInfo) error {
		if !fi.IsDir() {
			return nil
		}

		ver, err := ParseVersionPath(ev, path)
		if err != nil {
			// if the tag is not formatted as semver, skip it.
			return nil
		}
		if ver.Semver().GT(newest.Semver()) {
			newest = *ver
			found = true
		}
		return filepath.SkipDir
	}); err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrVersionNotFound
	}
	return &newest, nil
}

func FindVersion(ev Env, spec VersionSpec) (*Version, error) {
	if spec.tag != "" {
		return validateVersionSpec(ev, spec.owner, spec.name, spec.tag)
	}

	return findLatestVersion(ev, spec.owner, spec.name)
}

func ParseVersionPath(ev Env, path string) (*Version, error) {
	rel, err := filepath.Rel(assetSubPath(ev), path)
	if err != nil {
		return nil, err
	}

	terms := strings.Split(rel, string([]rune{filepath.Separator}))
	if len(terms) < 3 {
		return nil, errors.New("too short version path")
	}
	sv, err := ValidateTag(terms[2])
	if err != nil {
		return nil, errors.New("invalid version path")
	}
	ver := Version{
		App: App{
			owner: terms[0],
			name:  terms[1],
		},
		tag:    terms[2],
		semver: sv,
	}
	return &ver, nil
}

type VersionWalker func(Version) error

func WalkVersions(ev Env, walk VersionWalker) error {
	return walkIfDir(assetSubPath(ev), func(path string, fi os.FileInfo) error {
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
	return walkIfDir(AppPath(ev, app), func(path string, fi os.FileInfo) error {
		if !fi.IsDir() {
			return nil
		}
		tag := filepath.Base(path)
		if _, err := ValidateTag(tag); err != nil {
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
	uniq := map[string]struct{}{}
	if err := walkIfDir(ev.Bin(), walkLinkedVersions(ev, uniq, walk)); err != nil {
		return err
	}
	if err := walkIfDir(ev.Man(), walkLinkedVersions(ev, uniq, walk)); err != nil {
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
func walkLinkedVersions(ev Env, uniq map[string]struct{}, walk func(Version) error) func(path string, fi os.FileInfo) error {
	return func(path string, fi os.FileInfo) error {
		if (fi.Mode() & os.ModeSymlink) != os.ModeSymlink {
			return nil
		}
		link, err := os.Readlink(path)
		if err != nil {
			return err
		}

		ver, err := ParseVersionPath(ev, link)
		if err != nil {
			return nil // ignore invalid link
		}

		if _, ok := uniq[ver.String()]; ok {
			return nil
		}
		if err := walk(*ver); err != nil {
			return err
		}
		uniq[ver.String()] = struct{}{}
		return nil
	}
}
