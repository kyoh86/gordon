package command

import (
	"context"
	"os"

	"github.com/kyoh86/gordon/internal/gordon"
)

func Uninstall(ctx context.Context, ev Env, spec gordon.AppSpec) error {
	app, err := gordon.FindApp(ev, spec)
	if err != nil {
		return err
	}

	if err := gordon.Unlink(ev, *app); err != nil {
		return err
	}

	if err := os.RemoveAll(gordon.AppPath(ev, *app)); err != nil {
		return err
	}
	return nil
}
