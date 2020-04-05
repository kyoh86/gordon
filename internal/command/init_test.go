package command_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	keyring "github.com/zalando/go-keyring"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	keyring.MockInit()
	code := m.Run()
	os.Exit(code)
}
