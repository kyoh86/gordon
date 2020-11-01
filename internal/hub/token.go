package hub

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/99designs/keyring"
	"github.com/kyoh86/gordon/internal/cli"
	"github.com/kyoh86/xdg"
)

const (
	KeyringService = "gordon.kyoh86.dev"
	KeyringFileDir = "gordon"
)

func openKeyring(host string) (keyring.Keyring, error) {
	serviceName := strings.Join([]string{host, KeyringService}, ".")
	return keyring.Open(keyring.Config{
		ServiceName: serviceName,

		FileDir:          filepath.Join(xdg.CacheHome(), KeyringFileDir, "keyring", host),
		FilePasswordFunc: keyring.PromptFunc(cli.AskPassword),

		KeychainName:         serviceName,
		KeychainPasswordFunc: keyring.PromptFunc(cli.AskPassword),

		PassDir: filepath.Join(xdg.CacheHome(), KeyringFileDir, "pass", host),
	})
}

func SetGithubToken(host, user, token string) error {
	if host == "" {
		return errors.New("host is empty")
	}
	if user == "" {
		return errors.New("user is empty")
	}
	ring, err := openKeyring(host)
	if err != nil {
		return err
	}
	return ring.Set(keyring.Item{
		Key:  user,
		Data: []byte(token),
	})
}

func GetGithubToken(host, user string) (string, error) {
	if user == "" {
		return "", errors.New("user is empty")
	}
	envar := os.Getenv("GORDON_GITHUB_TOKEN")
	if envar != "" {
		return envar, nil
	}
	ring, err := openKeyring(host)
	if err != nil {
		return "", err
	}
	item, err := ring.Get(user)
	if err != nil {
		if errors.Is(err, keyring.ErrKeyNotFound) {
			return "", nil
		}
		return "", err
	}
	return string(item.Data), nil
}

func DeleteGithubToken(host, user string) error {
	ring, err := openKeyring(host)
	if err != nil {
		return err
	}
	return ring.Remove(user)
}
