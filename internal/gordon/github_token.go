package gordon

import (
	"errors"
	"os"
	"strings"

	"github.com/kyoh86/gordon/internal/env"
	"github.com/zalando/go-keyring"
)

func SetGitHubToken(host, user, token string) error {
	if host == "" {
		return errors.New("host is empty")
	}
	if user == "" {
		return errors.New("user is empty")
	}
	return keyring.Set(strings.Join([]string{host, env.KeyringService}, "."), user, token)
}

func GetGitHubToken(host, user string) (string, error) {
	if user == "" {
		return "", errors.New("user is empty")
	}
	envar := os.Getenv("GORDON_GITHUB_TOKEN")
	if envar != "" {
		return envar, nil
	}
	return keyring.Get(strings.Join([]string{host, env.KeyringService}, "."), user)
}
