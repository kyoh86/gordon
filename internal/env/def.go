package env

import (
	"os"
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

type Bin struct {
	extypes.Path
}

func (*Bin) Default() interface{} {
	return filepath.Join(os.Getenv("HOME"), ".local", "bin")
}

type Man struct {
	extypes.Path
}

func (*Man) Default() interface{} {
	return filepath.Join(os.Getenv("HOME"), ".local", "man")
}

type Cache struct {
	extypes.Path
}

func (*Cache) Default() interface{} {
	return filepath.Join(xdg.CacheHome(), "gordon")
}

type Hooks struct {
	extypes.Paths
}

func (*Hooks) Default() interface{} {
	return []string{filepath.Join(xdg.ConfigHome(), "gordon", "hooks")}
}
