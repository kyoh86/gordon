package command

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v29/github"
	"github.com/kyoh86/gordon/internal/gordon"
	"github.com/kyoh86/gordon/internal/hub"
)

func UpdateAll(ctx context.Context, ev Env) error {
	client, err := hub.NewClient(ctx, ev)
	if err != nil {
		return err
	}

	if err := gordon.WalkInstalledVersions(ev, func(ver gordon.Version) error {
		return update(ctx, client, ev, ver.Spec())
	}); err != nil {
		return err
	}
	return nil
}

func Update(ctx context.Context, ev Env, spec gordon.AppSpec) error {
	client, err := hub.NewClient(ctx, ev)
	if err != nil {
		return err
	}

	return update(ctx, client, ev, gordon.VersionSpec{
		AppSpec: spec,
	})
}

func update(ctx context.Context, client *github.Client, ev Env, spec gordon.VersionSpec) error {
	spec = spec.WithoutTag()
	release, err := gordon.FindRelease(ctx, ev, client, spec)
	if err != nil {
		return fmt.Errorf("failed to find release for %q-%q: %w", ev.OS(), ev.Architecture(), err)
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

	if !exist {
		log.Printf("info: %q has an old version %q\n", version.App, version.Tag())
		version, err = gordon.Download(ctx, ev, client, *release)
		if err != nil {
			return err
		}
		if err := gordon.Link(ev, *version); err != nil {
			return err
		}
		log.Printf("info: updated %q with new version %s\n", version.App, version.Tag())
	}

	return nil
}
