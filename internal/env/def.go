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

var (
	// DefaultHistoryFile is the default file to save history
	DefaultHistoryFile = filepath.Join(xdg.CacheHome(), "gordon", "history")
)

type HistoryFile struct{ types.StringPropertyBase }

func (*HistoryFile) Default() interface{} {
	return DefaultHistoryFile
}

type HistorySave struct{ types.BoolPropertyBase }

func (*HistorySave) Default() interface{} {
	return true
}

type ExtractModes struct{ FileModes }

func (*ExtractModes) Default() interface{} {
	return []os.FileMode{111}
}

type ExtractExclude struct{ types.StringPropertyBase }
type ExtractInclude struct{ types.StringPropertyBase }
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
