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
		if version.Tag() == release.Tag() {
			exist = true
		}
	case gordon.ErrVersionNotFound:
		// noop
	default:
		return err
	}

	if exist {
		log.Printf("info: %q has already installed with version %q\n", version.App, version.Tag())
	} else {
		version, err = gordon.Download(ctx, ev, client, *release)
		if err != nil {
			return err
		}
		if err := gordon.Link(ev, *version); err != nil {
			return err
		}
		log.Printf("info: installed %q with version %s\n", version.App, version.Tag())
	}

	return nil
}
