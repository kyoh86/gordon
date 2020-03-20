package command

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/google/go-github/v24/github"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gordon/internal/archive"
	"github.com/kyoh86/gordon/internal/context"
	"github.com/kyoh86/gordon/internal/gh"
	"github.com/kyoh86/gordon/internal/history"
)

// Download a package from GitHub Release.
// If `tag` is empty, it will download from the latest release.
func Download(ctx context.Context, repo *gogh.Repo, tag string, update bool) error {
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
	for _, asset := range release.Assets {
		if opener := assetOpener(ctx, asset); opener != nil {
			if err := download(ctx, repo, asset, opener, update); err != nil {
				return err
			}
			if err := history.SaveHistory(ctx, repo, tag); err != nil {
				log.Printf("warn: failed to save history: %v", err)
			}
			return nil
		}
	}
	return fmt.Errorf("there's no installable asset in release %s", release.GetTagName())
}

func assetOpener(ctx context.Context, asset github.ReleaseAsset) archive.Opener {
	name := asset.GetName()
	if !strings.Contains(name, ctx.Architecture()) {
		return nil
	}
	if !strings.Contains(name, ctx.OS()) {
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

func reg(pattern string) (*regexp.Regexp, error) {
	if pattern == "" {
		return nil, nil
	}
	return regexp.Compile(pattern)
}

var mkdirAllOnce sync.Once

func download(ctx context.Context, repo *gogh.Repo, asset github.ReleaseAsset, opener archive.Opener, update bool) error {
	log.Printf("info: download %s", asset.GetName())
	reader, err := gh.Asset(ctx, repo, asset.GetID())
	if err != nil {
		return err
	}
	defer reader.Close()

	arch, err := opener(reader)
	if err != nil {
		return err
	}

	excReg, err := reg(ctx.ExtractExclude())
	if err != nil {
		return err
	}
	incReg, err := reg(ctx.ExtractInclude())
	if err != nil {
		return err
	}
	return arch.Walk(func(info os.FileInfo, entry archive.Entry) (retErr error) {
		log.Printf("debug: extract %s", info.Name())
		if !ctx.ExtractModes().Match(info.Mode()) {
			log.Printf("debug: skip %s because mode %s is not matched", info.Name(), info.Mode())
			return nil
		}

		if excReg != nil && excReg.MatchString(info.Name()) {
			log.Printf("debug: exclude %s", info.Name())
			return nil
		}
		if incReg != nil && !incReg.MatchString(info.Name()) {
			log.Printf("debug: not included %s", info.Name())
			return nil
		}
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
			retErr = os.MkdirAll(ctx.Root(), 0777)
		})
		if retErr != nil {
			return
		}
		bin := filepath.Join(ctx.Root(), info.Name())
		flag := os.O_CREATE | os.O_EXCL | os.O_WRONLY
		if update {
			flag = os.O_TRUNC | os.O_WRONLY
		}
		file, err := os.OpenFile(bin, flag, info.Mode())
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
