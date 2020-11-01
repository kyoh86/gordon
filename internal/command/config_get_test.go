package command_test

import (
	"log"
	"strings"
	"testing"

	"github.com/kyoh86/gordon/internal/command"
	"github.com/kyoh86/gordon/internal/env"
	"github.com/stretchr/testify/assert"
)

func ExampleConfigGet() {
	yml := strings.NewReader(`{"hooks": ["/foo", "/bar"]}`)
	config, err := env.GetConfig(yml)
	if err != nil {
		log.Fatalln(err)
	}
	if err := command.ConfigGet(&config, "hooks"); err != nil {
		log.Fatalln(err)
	}

	// Output:
	// /foo:/bar
}

func TestConfigGet(t *testing.T) {
	config, err := env.GetConfig(env.EmptyYAMLReader)
	if err != nil {
		log.Fatalln(err)
	}
	assert.EqualError(t, command.ConfigGet(&config, "invalid.name"), `invalid option name "invalid.name"`)
}
