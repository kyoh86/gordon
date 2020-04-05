package command

import (
	"context"
	"log"

	"github.com/kyoh86/gordon/internal/gordon"
	"github.com/kyoh86/gordon/internal/hub"
)

func Install(ctx context.Context, ev Env, spec gordon.VersionSpec) error {
	client, err := hub.NewClient(ctx, ev)
	if err != nil {
		return err
	}

	release, err := gordon.FindRelease(ctx, ev, client, spec)
	if err != nil {
		return err
	}

	exist := false
	version, err := gordon.FindVersion(ev, spec)
	switch err {
	case nil:
		log.Printf("info: old version %q\n", version.Tag())
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
}