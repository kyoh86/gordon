package command_test

import (
	"testing"

	"github.com/kyoh86/gordon/internal/command"
	"github.com/kyoh86/gordon/internal/context"
	"github.com/stretchr/testify/assert"
)

func ExampleConfigGet() {
	if err := command.ConfigGet(&context.Config{
		VRoot: "/foo",
	}, "root"); err != nil {
		panic(err)
	}
	// Output:
	// /foo
}

func TestConfigGet(t *testing.T) {
	assert.EqualError(t, command.ConfigGet(&context.Config{}, "invalid.name"), "invalid option name")
}
