package command_test

import (
	"log"
	"strings"
	"testing"

	"github.com/kyoh86/gordon/internal/command"
	"github.com/kyoh86/gordon/internal/env"
	"github.com/stretchr/testify/assert"
)

func ExampleConfigGetAll() {
	yml := strings.NewReader(`
hooks:
  - /hook1
  - /hook2
githubHost: hostx1
githubUser: userx1`)
	config, err := env.GetConfig(yml)
	if err != nil {
		log.Fatalln(err)
	}
	if err := command.ConfigGetAll(nil, &config); err != nil {
		log.Fatalln(err)
	}

	// Unordered output:
	// hooks: /hook1:/hook2
	// github.host: hostx1
	// github.user: userx1
	// github.token: *****
	// architecture:
	// os:
	// cache:
	// bin:
	// man:
}

func TestConfigGetAll(t *testing.T) {
	yml := strings.NewReader(`
hooks:
  - /hook1
  - /hook2
githubHost: hostx1`)
	config, err := env.GetConfig(yml)
	assert.NoError(t, err)
	assert.NoError(t, command.ConfigGetAll(nil, &config))
}
