package env

import (
	"path/filepath"
	"runtime"

	"github.com/kyoh86/appenv/extypes"
	"github.com/kyoh86/appenv/types"
	"github.com/kyoh86/xdg"
)

const (
	KeyringService = "gordon.kyoh86.dev"
	EnvarPrefix    = "GORDON_"
)

type GithubHost struct {
	types.StringPropertyBase
}

const (
	// DefaultHost is the default host of the GitHub
	DefaultHost = "github.com"
)

func (*GithubHost) Default() interface{} {
	return DefaultHost
}

type GithubUser struct {
	types.StringPropertyBase
}

type Architecture struct{ types.StringPropertyBase }

func (*Architecture) Default() interface{} {
	return runtime.GOARCH
}

type OS struct{ types.StringPropertyBase }

func (*OS) Default() interface{} {
	return runtime.GOOS
}

type Root struct {
	extypes.Path
}

func (*Root) Default() interface{} {
	return xdg.DownloadDir()
}

type Hooks struct {
	extypes.Paths
}

func (*Hooks) Default() interface{} {
	return []string{filepath.Join(xdg.ConfigHome(), "gordon", "hooks")}
}
