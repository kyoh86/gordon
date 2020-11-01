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

func ExampleConfigUnset() {
	source := strings.NewReader(`
hooks:
  - /hook1
githubUser: userx1
githubHost: hostx1`)
	config, access, err := env.GetAppenv(source, env.EnvarPrefix)
	if err != nil {
		log.Fatalln(err)
	}
	if err := command.ConfigUnset(&access, &config, "github.host"); err != nil {
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
	// githubUser: userx1
	// hooks: /hook1
	// github.host:
	// github.user: userx1
	// github.token: *****
	// architecture:
	// os:
	// cache:
	// bin:
	// man:
}

func TestConfigUnset(t *testing.T) {
	source := strings.NewReader(`
hooks:
    - /hook1
githubUser: userx1
githubHost: hostx1`)
	config, access, err := env.GetAppenv(source, env.EnvarPrefix)
	assert.NoError(t, err)
	assert.NoError(t, command.ConfigUnset(&access, &config, "github.host"))
	assert.NoError(t, config.Save(os.Stdout))
	assert.NoError(t, command.ConfigGetAll(nil, &config))
	assert.Error(t, command.ConfigUnset(&access, &config, "invalid.config"))
}
