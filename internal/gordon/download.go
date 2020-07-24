package gordon

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/v29/github"
	"github.com/kyoh86/gordon/internal/archive"
)

func Download(ctx context.Context, ev Env, client *github.Client, release Release) (*Version, error) {
	version := ReleaseVersion(release)
	path := VersionPath(ev, version)
	// UNDONE: if the version already exist
	if err := os.MkdirAll(path, 0777); err != nil {
		return nil, err
	}
	unarchiver, err := openAsset(ctx, client, release)
	if err != nil {
		return nil, fmt.Errorf("open asset %s: %w", release.Spec().String(), err)
	}
	if err := extractAsset(path, unarchiver); err != nil {
		return nil, err
	}
	return &version, nil
}

func openAsset(ctx context.Context, client *github.Client, release Release) (archive.Unarchiver, error) {
	reader, _, err := client.Repositories.DownloadReleaseAsset(
		ctx,
		release.owner,
		release.name,
		release.asset.ID,
		http.DefaultClient,
	)
	if err != nil {
		return nil, err
	}
	return release.opener(reader)
}

func extractAsset(path string, unarchiver archive.Unarchiver) error {
	if err := os.MkdirAll(path, 0777); err != nil {
		return err
	}

	if err := unarchiver.Walk(func(info os.FileInfo, entry archive.Entry) (retErr error) {
		entryReader, err := entry()
		if err != nil {
			return fmt.Errorf("open an entry: %w", err)
		}
		defer func() {
			if err := entryReader.Close(); err != nil && retErr == nil {
				retErr = fmt.Errorf("close an entry: %w", err)
			}
		}()

		// TODO: support subdirectory
		file, err := os.OpenFile(
			filepath.Join(path, info.Name()),
			os.O_CREATE|os.O_EXCL|os.O_WRONLY,
			info.Mode(),
		)
		if err != nil {
			return fmt.Errorf("open a destination: %w", err)
		}
		defer func() {
			if err := file.Close(); err != nil && retErr == nil {
				retErr = fmt.Errorf("close a destination: %w", err)
			}
		}()

		if _, err := io.Copy(file, entryReader); err != nil {
			return fmt.Errorf("copy from an entry: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("unarchive %s: %w", path, err)
	}
	return nil
}
