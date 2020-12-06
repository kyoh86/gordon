package gordon

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-github/v29/github"
	"github.com/kyoh86/gordon/internal/archive"
)

type Release struct {
	Repo
	tag string

	asset  Asset
	opener archive.Opener
}

func (r Release) String() string {
	return fmt.Sprintf("%s/%s", r.Spec().String(), r.asset.Name)
}

func (r Release) Spec() VersionSpec {
	return VersionSpec{
		AppSpec: r.Repo.Spec(),
		tag:     r.tag,
	}
}

func ReleaseVersion(release Release) Version {
	return Version{
		App: RepoApp(release.Repo),
		tag: release.tag,
	}
}

type Asset struct {
	ID                 int64
	Name               string
	Label              string
	ContentType        string
	Size               int
	BrowserDownloadURL string
}

func (a Release) Owner() string { return a.owner }
func (a Release) Name() string  { return a.name }
func (a Release) Tag() string   { return a.tag }

var (
	ErrAssetNotFound = errors.New("asset not found")
)

func FindRelease(ctx context.Context, ev Env, client *github.Client, spec VersionSpec) (*Release, error) {
	release, err := findRelease(ctx, client, spec)
	if err != nil {
		return nil, err
	}

	asset, err := findAsset(ev, release.Assets)
	if err != nil {
		return nil, err
	}
	opener := matchOpener(asset.GetName())
	return &Release{
		Repo: Repo{
			owner: spec.owner,
			name:  spec.name,
		},
		tag: release.GetTagName(),

		asset: Asset{
			ID:                 asset.GetID(),
			Name:               asset.GetName(),
			Label:              asset.GetLabel(),
			ContentType:        asset.GetContentType(),
			Size:               asset.GetSize(),
			BrowserDownloadURL: asset.GetBrowserDownloadURL(),
		},
		opener: opener,
	}, nil
}

func findRelease(ctx context.Context, client *github.Client, spec VersionSpec) (*github.RepositoryRelease, error) {
	if spec.tag == "" {
		//UNDONE: list releases and find newest vertag
		release, _, err := client.Repositories.GetLatestRelease(ctx, spec.Owner(), spec.Name())
		return release, err
	}
	release, _, err := client.Repositories.GetReleaseByTag(ctx, spec.Owner(), spec.Name(), spec.Tag())
	return release, err
}
