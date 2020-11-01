package command_test

import (
	"os"
	"testing"

	testtarget "github.com/kyoh86/gordon/internal/command"
	"github.com/kyoh86/gordon/internal/hub"
)

func TestMain(m *testing.M) {
	testtarget.TokenManager = hub.NewMemory
	code := m.Run()
	os.Exit(code)
}
