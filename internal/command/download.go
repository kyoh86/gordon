package command

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/go-github/v29/github"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gordon/internal/archive"
	"github.com/kyoh86/gordon/internal/env"
	"github.com/kyoh86/gordon/internal/hub"
)

// Download a package from GitHub Release.
// If `tag` is empty, it will download from the latest release.
func Download(ctx context.Context, ev env.Env, spec gogh.RepoSpec, tag string, update bool) error {
	repo, err := spec.Validate(ev)
	if err != nil {
		return err
	}
	var release *github.RepositoryRelease
	client, err := hub.New(ctx, ev)
	if err != nil {
		return err
	}
	if tag == "" {
		rel, err := client.LatestRelease(ctx, repo)
		if err != nil {
			return err
		}
		release = rel
	} else {
		rel, err := client.Release(ctx, repo, tag)
		if err != nil {
			return err
		}
		release = rel
	}
	for _, asset := range release.Assets {
		if opener := assetOpener(ev, asset); opener != nil {
			if err := download(ctx, ev, client, repo, asset, opener, update); err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("there's no installable asset in release %s", release.GetTagName())
}

func assetOpener(ev env.Env, asset github.ReleaseAsset) archive.Opener {
	name := asset.GetName()
	if !strings.Contains(name, ev.Architecture()) {
		return nil
	}
	if !strings.Contains(name, ev.OS()) {
		return nil
	}

	switch {
	case strings.HasSuffix(name, ".tar.gz"), strings.HasSuffix(name, ".tgz"):
		return archive.OpenTarGzip
	case strings.HasSuffix(name, ".tar.bz2"):
		return archive.OpenTarBzip2
	case strings.HasSuffix(name, ".tar"):
		return archive.OpenTar
	case strings.HasSuffix(name, ".zip"):
		return archive.ZipOpener(int64(asset.GetSize()))
	}
	return nil
}

var mkdirAllOnce sync.Once

func download(ctx context.Context, ev env.Env, client *hub.Client, repo *gogh.Repo, asset github.ReleaseAsset, opener archive.Opener, update bool) error {
	log.Printf("info: download %s", asset.GetName())
	reader, err := client.Asset(ctx, repo, asset.GetID())
	if err != nil {
		return err
	}
	defer reader.Close()

	arch, err := opener(reader)
	if err != nil {
		return err
	}

	return arch.Walk(func(info os.FileInfo, entry archive.Entry) (retErr error) {
		log.Printf("info: unarchive %s", info.Name())

		entryReader, err := entry()
		if err != nil {
			return err
		}
		defer func() {
			if err := entryReader.Close(); err != nil && retErr == nil {
				retErr = err
			}
		}()
		mkdirAllOnce.Do(func() {
			retErr = os.MkdirAll(ev.Root(), 0777)
		})
		if retErr != nil {
			return
		}
		path := filepath.Join(ev.Root(), info.Name())
		flag := os.O_CREATE | os.O_EXCL | os.O_WRONLY
		if update {
			flag = os.O_TRUNC | os.O_WRONLY
		}

		file, err := os.OpenFile(path, flag, info.Mode())
		if err != nil {
			return err
		}
		defer func() {
			if err := file.Close(); err != nil && retErr == nil {
				retErr = err
			}
		}()

		if _, err := io.Copy(file, entryReader); err != nil {
			return err
		}
		return nil
	})
}
