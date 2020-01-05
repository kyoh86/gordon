package context

import (
	"context"
	"io"
	"runtime"

	"github.com/joeshaw/envdecode"
	"github.com/kyoh86/xdg"
	yaml "gopkg.in/yaml.v2"
)

var (
	envGordonLogLevel        = "GORDON_LOG_LEVEL"
	envGordonLogDate         = "GORDON_LOG_DATE"
	envGordonLogTime         = "GORDON_LOG_TIME"
	envGordonLogMicroSeconds = "GORDON_LOG_MICROSECONDS"
	envGordonLogLongFile     = "GORDON_LOG_LONGFILE"
	envGordonLogShortFile    = "GORDON_LOG_SHORTFILE"
	envGordonLogUTC          = "GORDON_LOG_UTC"
	envGordonGitHubUser      = "GORDON_GITHUB_USER"
	envGordonGitHubToken     = "GORDON_GITHUB_TOKEN"
	envGordonGitHubHost      = "GORDON_GITHUB_HOST"
	envGordonRoot            = "GORDON_ROOT"
	envGordonArch            = "GORDON_ARCH"
	envGordonOS              = "GORDON_OS"
	envNames                 = []string{
		envGordonLogLevel,
		envGordonLogDate,
		envGordonLogTime,
		envGordonLogMicroSeconds,
		envGordonLogLongFile,
		envGordonLogShortFile,
		envGordonLogUTC,
		envGordonGitHubUser,
		envGordonGitHubToken,
		envGordonGitHubHost,
		envGordonRoot,
		envGordonArch,
		envGordonOS,
	}
)

const (
	// DefaultHost is the default host of the GitHub
	DefaultHost = "github.com"
	// DefaultLogLevel is the default level to output log
	DefaultLogLevel = "warn"
)

var defaultConfig = Config{
	Context: context.Background(),
	Log: LogConfig{
		Level: DefaultLogLevel,
		Time:  TrueOption,
	},
	GitHub: GitHubConfig{
		Host: DefaultHost,
	},
	VRoot: xdg.DownloadDir(),
	VArch: runtime.GOARCH,
	VOS:   runtime.GOOS,
}

func DefaultConfig() *Config {
	return &defaultConfig
}

func LoadConfig(r io.Reader) (config *Config, err error) {
	config = &Config{}
	if err := yaml.NewDecoder(r).Decode(config); err != nil {
		return nil, err
	}
	return
}

func SaveConfig(w io.Writer, config *Config) error {
	return yaml.NewEncoder(w).Encode(config)
}

func GetEnvarConfig() (config *Config, err error) {
	config = &Config{}
	err = envdecode.Decode(config)
	if err == envdecode.ErrNoTargetFieldsAreSet {
		err = nil
	}
	return
}
