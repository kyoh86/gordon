package hub

import (
	"context"
	"io"
	"net/http"

	"github.com/google/go-github/v29/github"
	"github.com/kyoh86/gogh/gogh"
)

// LatestRelease gets the latest release.
// Parameters:
//   * repo: Target repository.
func (c *Client) LatestRelease(ctx context.Context, repo *gogh.Repo) (*github.RepositoryRelease, error) {
	release, _, err := c.client.Repositories.GetLatestRelease(ctx, repo.Owner(), repo.Name())
	if err != nil {
		return nil, err
	}

	return release, nil
}

// Release gets the tagged release.
// Parameters:
//   * repo: Target repository.
//   * tag:  Target tag.
func (c *Client) Release(ctx context.Context, repo *gogh.Repo, tag string) (*github.RepositoryRelease, error) {
	release, _, err := c.client.Repositories.GetReleaseByTag(ctx, repo.Owner(), repo.Name(), tag)
	if err != nil {
		return nil, err
	}
	return release, nil
}

// Asset downloads an asset.
// Parameters:
//   * repo:    Target repository.
//   * assetID: Target asset ID.
func (c *Client) Asset(ctx context.Context, repo *gogh.Repo, assetID int64) (io.ReadCloser, error) {
	reader, _, err := c.client.Repositories.DownloadReleaseAsset(ctx, repo.Owner(), repo.Name(), assetID, http.DefaultClient)
	return reader, err
}
