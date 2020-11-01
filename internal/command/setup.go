package command

import (
	"context"
	"fmt"

	"github.com/kyoh86/ask"
	"github.com/kyoh86/gordon/internal/cli"
	"github.com/kyoh86/gordon/internal/env"
)

func Setup(_ context.Context, ev Env, cfg *env.Config, force bool) error {
	user := ev.GithubUser()
	if user == "" || force {
		if err := ask.Default(ev.GithubUser()).Message(cli.Question("Enter your GitHub user ID")).StringVar(&user).Do(); err != nil {
			return fmt.Errorf("asking GitHub user ID: %w", err)
		}

		return cfg.GithubUser().Set(user)
	}
	tm, err := TokenManager(ev.GithubHost())
	if err != nil {
		return err
	}
	token, _ := tm.GetGithubToken(user)
	if token == "" || force {
		newToken, err := cli.AskPassword("Enter your GitHub Private Access Token")
		if err != nil {
			return fmt.Errorf("asking GitHub Private Access Token: %w", err)
		}
		token = newToken
	}

	return tm.SetGithubToken(user, token)
}
