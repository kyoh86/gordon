package gh

import (
	"io"
	"net/http"

	"github.com/google/go-github/v24/github"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gordon/internal/context"
)

// LatestRelease gets the latest release.
// Parameters:
//   * repo: Target repository.
func LatestRelease(ctx context.Context, repo *gogh.Repo) (*github.RepositoryRelease, error) {
	client, err := NewClient(ctx)
	if err != nil {
		return nil, err
	}

	release, _, err := client.Repositories.GetLatestRelease(ctx, repo.Owner(ctx), repo.Name(ctx))
	if err != nil {
		return nil, err
	}

	// TODO: investigate that "Assets" is enough or not.
	// if "Assets" is not enough, use "List Assets for a release" API
	return release, nil
}

// Release gets the tagged release.
// Parameters:
//   * repo: Target repository.
//   * tag:  Target tag.
func Release(ctx context.Context, repo *gogh.Repo, tag string) (*github.RepositoryRelease, error) {
	client, err := NewClient(ctx)
	if err != nil {
		return nil, err
	}

	release, _, err := client.Repositories.GetReleaseByTag(ctx, repo.Owner(ctx), repo.Name(ctx), tag)
	if err != nil {
		return nil, err
	}
	return release, nil
}

// Asset downloads an asset.
// Parameters:
//   * repo:    Target repository.
//   * assetID: Target asset ID.
func Asset(ctx context.Context, repo *gogh.Repo, assetID int64) (io.ReadCloser, error) {
	client, err := NewClient(ctx)
	if err != nil {
		return nil, err
	}
	reader, redirect, err := client.Repositories.DownloadReleaseAsset(ctx, repo.Owner(ctx), repo.Name(ctx), assetID)
	if reader != nil {
		return reader, err
	}
	res, err := http.Get(redirect)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}
