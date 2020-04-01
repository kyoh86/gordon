package env

import (
	"go/build"
	"path/filepath"

	"github.com/kyoh86/appenv/extypes"
	"github.com/kyoh86/appenv/types"
	"github.com/kyoh86/xdg"
	"github.com/thoas/go-funk"
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

type Root struct {
	extypes.Path
}

func (*Root) Default() interface{} {
	gopaths := filepath.SplitList(build.Default.GOPATH)
	paths := make([]string, 0, len(gopaths))
	for _, gopath := range gopaths {
		paths = append(paths, filepath.Join(gopath, "src"))
	}
	return funk.UniqString(paths)
}

type Hooks struct {
	extypes.Paths
}

func (*Hooks) Default() interface{} {
	return []string{filepath.Join(xdg.ConfigHome(), "gordon", "hooks")}
}
