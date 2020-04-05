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
	"github.com/kyoh86/gordon/internal/gordon"
	keyring "github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

// NewClient builds GitHub Client with GitHub API token that is configured.
func NewClient(authContext context.Context, ev gordon.Env) (*github.Client, error) {
	if host := ev.GithubHost(); host != "" && host != "github.com" {
		url := fmt.Sprintf("https://%s/api/v3", host)
		httpClient, err := oauth2Client(authContext, ev)
		if err != nil {
			return nil, err
		}
		return github.NewEnterpriseClient(url, url, httpClient)
	}
	httpClient, err := oauth2Client(authContext, ev)
	if err != nil {
		return nil, err
	}
	return github.NewClient(httpClient), nil
}

// New builds hub.Client with GitHub API token that is configured.
func New(authContext context.Context, ev gordon.Env) (*Client, error) {
	client, err := NewClient(authContext, ev)
	if err != nil {
		return nil, err
	}
	return &Client{client}, nil
}

func getToken(ev gordon.Env) (string, error) {
	if ev.GithubUser() == "" {
		return "", errors.New("github.user is empty")
	}
	envar := os.Getenv("GORDON_GITHUB_TOKEN")
	if envar != "" {
		return envar, nil
	}
	return keyring.Get(strings.Join([]string{ev.GithubHost(), env.KeyringService}, "."), ev.GithubUser())
}

func oauth2Client(authContext context.Context, ev gordon.Env) (*http.Client, error) {
	if ev.GithubUser() == "" {
		return http.DefaultClient, nil
	}
	token, err := getToken(ev)
	if err != nil {
		return nil, err
	}
	if token == "" {
		return http.DefaultClient, nil
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	return oauth2.NewClient(authContext, ts), nil
}

type Client struct {
	client *github.Client
}
