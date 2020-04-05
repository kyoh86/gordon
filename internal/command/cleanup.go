package command

import (
	"os"

	"github.com/kyoh86/gordon/internal/gordon"
)

func Cleanup(ev Env) error {
	linked := map[string]struct{}{}
	if err := gordon.WalkInstalledVersions(ev, func(ver gordon.Version) error {
		linked[ver.String()] = struct{}{}
		return nil
	}); err != nil {
		return err
	}
	if err := gordon.WalkVersions(ev, func(ver gordon.Version) error {
		if _, newest := linked[ver.String()]; newest {
			return nil
		}
		return os.RemoveAll(gordon.VersionPath(ev, ver))
	}); err != nil {
		return err
	}
	return nil
}
