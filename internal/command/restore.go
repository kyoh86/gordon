package command

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"

	"github.com/kyoh86/gordon/internal/gordon"
	"github.com/kyoh86/gordon/internal/hub"
)

func Restore(ctx context.Context, ev Env, bundle string) error {
	client, err := hub.NewClient(ctx, ev)
	if err != nil {
		return err
	}
	var reader io.Reader
	if bundle == "-" {
		reader = os.Stdin
	} else {
		f, err := os.Open(bundle)
		if err != nil {
			return err
		}
		reader = f
		defer f.Close()
	}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		spec, err := gordon.ParseVersionSpec(scanner.Text())
		if err != nil {
			return err
		}
		release, err := gordon.FindRelease(ctx, ev, client, *spec)
		if err != nil {
			return err
		}

		version, err := gordon.Download(ctx, ev, client, *release)
		if err != nil {
			return err
		}

		if err := gordon.Link(ev, *version); err != nil {
			return err
		}
		log.Printf("info: restored %q with version %s\n", version.App, version.Tag())
	}
	return nil
}
