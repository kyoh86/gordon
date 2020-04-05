package gordon

import (
	"context"
	"errors"
	"strings"

	"github.com/google/go-github/v29/github"
	"github.com/kyoh86/gordon/internal/archive"
)

type Release struct {
	Repo
	tag string

	asset  Asset
	opener archive.Opener
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

	for _, asset := range release.Assets {
		if opener := assetOpener(ev, asset); opener != nil {
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
	}
	return nil, ErrAssetNotFound
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

func assetOpener(ev Env, asset github.ReleaseAsset) archive.Opener {
	name := asset.GetName()
	if !MatchArchitecture(name, ev.Architecture()) {
		return nil
	}
	if !MatchOS(name, ev.OS()) {
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
