package command

import (
	"context"

	"github.com/kyoh86/gordon/internal/env"
)

func Install(ctx context.Context, ev env.Env, spec gordon.AppSpec) error {
	release, err := gordon.FindRelease(spec)
	if err != nil {
		return err
	}

	asset, err := gordon.FindAsset(ev, release)
	if err != nil {
		return err
	}

	storePath := release.GetStorePath(ev)
	if err := gordon.Download(ev, asset, storePath); err != nil {
		return err
	}
	if err := gordon.Link(ev, storePath); err != nil {
		return err
	}
	return nil
}
