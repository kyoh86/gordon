package command

import (
	"github.com/kyoh86/gordon/internal/gordon"
)

func Cleanup(ev Env) error {
	linked := map[gordon.Version]struct{}{}
	if err := gordon.WalkInstalledVersions(ev, func(ver gordon.Version) error {
		linked[ver] = struct{}{}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
