package command

import (
	"context"
	"log"

	"github.com/kyoh86/gordon/internal/gordon"
	"github.com/kyoh86/gordon/internal/hub"
)

func Update(ctx context.Context, ev Env) error {
	client, err := hub.NewClient(ctx, ev)
	if err != nil {
		return err
	}

	if err := gordon.WalkInstalledVersions(ev, func(ver gordon.Version) error {
		spec := ver.Spec().WithoutTag()
		release, err := gordon.FindRelease(ctx, ev, client, spec)
		if err != nil {
			return err
		}

		exist := false
		version, err := gordon.FindVersion(ev, spec)
		switch err {
		case nil:
			log.Printf("info: %q has a version %q\n", version.App, version.Tag())
			if version.Tag() == release.Tag() {
				exist = true
			}
		case gordon.ErrVersionNotFound:
			// noop
		default:
			return err
		}

		if !exist {
			version, err = gordon.Download(ctx, ev, client, *release)
			if err != nil {
				return err
			}
		}

		if err := gordon.Link(ev, *version); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
