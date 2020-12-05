package command

import (
	"context"
	"fmt"

	"github.com/kyoh86/gordon/internal/gordon"
	"github.com/kyoh86/gordon/internal/hub"
)

func Get(ctx context.Context, ev Env, spec gordon.VersionSpec) error {
	client, err := hub.NewClient(ctx, ev)
	if err != nil {
		return err
	}

	release, err := gordon.FindRelease(ctx, ev, client, spec)
	if err != nil {
		return fmt.Errorf("failed to find release for %q-%q: %w", ev.OS(), ev.Architecture(), err)
	}

	ver, err := gordon.Download(ctx, ev, client, *release)
	if err != nil {
		return err
	}
	_ = ver
	return nil
}
