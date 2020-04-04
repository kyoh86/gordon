package command

import (
	"context"

	"github.com/kyoh86/gordon/internal/gordon"
	"github.com/kyoh86/gordon/internal/hub"
)

func Install(ctx context.Context, ev Env, spec gordon.AppSpec) error {
	client, err := hub.NewClient(ctx, ev)
	if err != nil {
		return err
	}

	release, err := gordon.FindRelease(ctx, ev, client, spec)
	if err != nil {
		return err
	}

	if err := gordon.Download(ctx, ev, client, release); err != nil {
		return err
	}
	if err := gordon.Link(ev, release); err != nil {
		return err
	}
	return nil
}
