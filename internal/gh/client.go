package gh

import (
	"fmt"
	"net/http"

	"github.com/google/go-github/v24/github"
	"github.com/kyoh86/gordon/internal/context"
	"golang.org/x/oauth2"
)

// NewClient builds GitHub Client with GitHub API token that is configured.
func NewClient(ctx context.Context) (*github.Client, error) {
	if host := ctx.GitHubHost(); host != "" && host != "github.com" {
		url := fmt.Sprintf("https://%s/api/v3", host)
		return github.NewEnterpriseClient(url, url, oauth2Client(ctx))
	}

	return github.NewClient(oauth2Client(ctx)), nil
}

func authenticated(ctx context.Context) bool {
	return ctx.GitHubToken() != ""
}

func oauth2Client(ctx context.Context) *http.Client {
	if !authenticated(ctx) {
		return nil
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ctx.GitHubToken()})
	return oauth2.NewClient(ctx, ts)
}
