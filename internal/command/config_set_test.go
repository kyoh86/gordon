package command_test

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/kyoh86/gordon/internal/command"
	"github.com/kyoh86/gordon/internal/env"
	"github.com/stretchr/testify/assert"
)

func ExampleConfigSet() {
	source := strings.NewReader(`
hooks:
  - /hook1
githubUser: userx1
githubHost: hostx1`)
	config, access, err := env.GetAppenv(source, env.EnvarPrefix)
	if err != nil {
		log.Fatalln(err)
	}
	if err := command.ConfigSet(&access, &config, "github.host", "hostx2"); err != nil {
		log.Fatalln(err)
	}
	if err := config.Save(os.Stdout); err != nil {
		log.Fatalln(err)
	}
	if err := command.ConfigGetAll(nil, &config); err != nil {
		log.Fatalln(err)
	}

	// Unordered output:
	// hooks:
	// - /hook1
	// githubHost: hostx2
	// githubUser: userx1
	// hooks: /hook1
	// github.host: hostx2
	// github.user: userx1
	// github.token: *****
	// architecture:
	// os:
	// cache:
	// bin:
	// man:
}

func TestConfigSet(t *testing.T) {
	// NOTE: never use real host name. github.token breaks keyring store
	source := strings.NewReader(`
hooks:
  - /hook1
githubUser: userx1
githubHost: hostx1`)
	config, access, err := env.GetAppenv(source, env.EnvarPrefix)
	assert.NoError(t, err)
	assert.NoError(t, command.ConfigSet(&access, &config, "github.host", "hostx2"))
	assert.NoError(t, config.Save(os.Stdout))
	assert.NoError(t, command.ConfigGetAll(nil, &config))

	assert.Error(t, command.ConfigSet(&access, &config, "invalid.config", "invalid"))
	assert.NoError(t, command.ConfigSet(&access, &config, "github.token", "invalid"))
	assert.NoError(t, command.ConfigUnset(&access, &config, "github.token"))
}
