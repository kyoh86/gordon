package hub

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/v29/github"
	"github.com/kyoh86/gordon/internal/env"
	keyring "github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

// New builds GitHub Client with GitHub API token that is configured.
func New(authContext context.Context, ev env.Env) (*Client, error) {
	if host := ev.GithubHost(); host != "" && host != "github.com" {
		url := fmt.Sprintf("https://%s/api/v3", host)
		httpClient, err := oauth2Client(authContext, ev)
		if err != nil {
			return nil, err
		}
		client, err := github.NewEnterpriseClient(url, url, httpClient)
		if err != nil {
			return nil, err
		}
		return &Client{client}, nil
	}
	httpClient, err := oauth2Client(authContext, ev)
	if err != nil {
		return nil, err
	}
	return &Client{github.NewClient(httpClient)}, nil
}

func getToken(ev env.Env) (string, error) {
	if ev.GithubUser() == "" {
		return "", errors.New("github.user is empty")
	}
	envar := os.Getenv("GOGH_GITHUB_TOKEN")
	if envar != "" {
		return envar, nil
	}
	return keyring.Get(strings.Join([]string{ev.GithubHost(), env.KeyringService}, "."), ev.GithubUser())
}

func oauth2Client(authContext context.Context, ev env.Env) (*http.Client, error) {
	token, err := getToken(ev)
	if err != nil {
		return nil, err
	}
	if token == "" {
		return nil, errors.New("github.token is empty")
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	return oauth2.NewClient(authContext, ts), nil
}

type Client struct {
	client *github.Client
}
