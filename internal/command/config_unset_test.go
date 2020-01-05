package command_test

import (
	"testing"

	"github.com/kyoh86/gordon/internal/command"
	"github.com/kyoh86/gordon/internal/context"
	"github.com/stretchr/testify/assert"
)

func TestConfigUnset(t *testing.T) {
	cfg := context.Config{
		GitHub: context.GitHubConfig{
			Host: "hostx1",
		},
	}
	assert.NoError(t, command.ConfigUnset(&cfg, "github.host"))
	assert.Empty(t, cfg.GitHub.Host)
	assert.EqualError(t, command.ConfigUnset(&cfg, "invalid.name"), "invalid option name")
}
