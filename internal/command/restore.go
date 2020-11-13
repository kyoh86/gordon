package command

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"sync"

	"github.com/google/go-github/v29/github"
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
	mute := sync.Mutex{}
	uniq := map[string]struct{}{}
	if err := gordon.WalkVersions(ev, func(v gordon.Version) error {
		mute.Lock()
		defer mute.Unlock()
		uniq[v.String()] = struct{}{}
		return nil
	}); err != nil {
		return err
	}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if _, exist := uniq[scanner.Text()]; exist {
			continue
		}
		if err := restoreOne(ctx, ev, client, scanner.Text()); err != nil {
			log.Printf("error: failed to restore %s: %s", scanner.Text(), err)
		}
	}
	return nil
}

func restoreOne(ctx context.Context, ev Env, client *github.Client, specStr string) error {
	spec, err := gordon.ParseVersionSpec(specStr)
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
	return nil
}
