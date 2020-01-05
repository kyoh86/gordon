package command

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v24/github"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gordon/internal/archive"
	"github.com/kyoh86/gordon/internal/context"
	"github.com/kyoh86/gordon/internal/gh"
)

// Download a package from GitHub Release.
// If `tag` is empty, it will download from the latest release.
func Download(ctx context.Context, repo *gogh.Repo, tag string) error {
	var release *github.RepositoryRelease
	if tag == "" {
		rel, err := gh.LatestRelease(ctx, repo)
		if err != nil {
			return err
		}
		release = rel
	} else {
		rel, err := gh.Release(ctx, repo, tag)
		if err != nil {
			return err
		}
		release = rel
	}
	//TODO: store download history with options
	//TODO: accept exclusion (asset name | tag) pattern
	//TODO: accept inclusion (asset name) pattern (e.g. *.ttf)
	//TODO: set download-root temporarily
	for _, asset := range release.Assets {
		if starter := assetOpener(ctx, asset); starter != nil {
			return download(ctx, repo, asset, starter)
		}
	}
	return fmt.Errorf("there's no installable asset in release %s", release.GetTagName())
}

func assetOpener(ctx context.Context, asset github.ReleaseAsset) archive.Opener {
	if !strings.Contains(asset.GetName(), ctx.Arch()) {
		return nil
	}
	if !strings.Contains(asset.GetName(), ctx.OS()) {
		return nil
	}
	name := asset.GetName()
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

func download(ctx context.Context, repo *gogh.Repo, asset github.ReleaseAsset, starter archive.Opener) error {
	reader, err := gh.Asset(ctx, repo, asset.GetID())
	if err != nil {
		return err
	}
	defer reader.Close()

	arch, err := starter(reader)
	if err != nil {
		return err
	}
	return arch.Walk(func(info os.FileInfo, entry archive.Entry) error {
		if (info.Mode() & 0111) == 0 {
			return nil
		}
		reader, err := entry()
		if err != nil {
			return err
		}
		defer reader.Close()
		bin := filepath.Join(ctx.Root(), info.Name())
		file, err := os.OpenFile(bin, os.O_CREATE|os.O_EXCL|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := io.Copy(file, reader); err != nil {
			return err
		}
		return nil
	})
}
