package context

import (
	"context"
	"io"
	"log"
	"path/filepath"
	"runtime"

	"github.com/joeshaw/envdecode"
	"github.com/kyoh86/xdg"
	"github.com/zalando/go-keyring"
	yaml "gopkg.in/yaml.v3"
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
	envGordonHistoryFile     = "GORDON_HISTORY_FILE"
	envGordonHistorySave     = "GORDON_HISTORY_SAVE"
	envGordonExtractModes    = "GORDON_EXTRACT_MODES"
	envGordonExtractExclude  = "GORDON_EXTRACT_EXCLUDE"
	envGordonExtractInclude  = "GORDON_EXTRACT_INCLUDE"
	envGordonRoot            = "GORDON_ROOT"
	envGordonArchitecture    = "GORDON_ARCHITECTURE"
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
		envGordonHistoryFile,
		envGordonHistorySave,
		envGordonExtractModes,
		envGordonExtractExclude,
		envGordonExtractInclude,
		envGordonRoot,
		envGordonArchitecture,
		envGordonOS,
	}

	keyGordonServiceName = "gordon.kyoh86.dev"
	keyGordonGitHubToken = "github-token"
)

const (
	// DefaultLogLevel is the default level to output log
	DefaultLogLevel = "warn"
	// DefaultGitHubHost is the default host of the GitHub
	DefaultGitHubHost = "github.com"
)

var (
	// DefaultHistoryFile is the default file to save history
	DefaultHistoryFile = filepath.Join(xdg.CacheHome(), "gordon", "history")
)

var defaultConfig = Config{
	Context: context.Background(),
	Log: LogConfig{
		Level: DefaultLogLevel,
		Time:  TrueOption,
	},
	GitHub: GitHubConfig{
		Host: DefaultGitHubHost,
	},
	History: HistoryConfig{
		File: DefaultHistoryFile,
		Save: TrueOption,
	},
	Extract: ExtractConfig{
		Modes: FileModes{0111},
	},
	VRoot:         xdg.DownloadDir(),
	VArchitecture: runtime.GOARCH,
	VOS:           runtime.GOOS,
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

func LoadKeyring() *Config {
	token, err := keyring.Get(keyGordonServiceName, keyGordonGitHubToken)
	if err != nil {
		log.Printf("info: there's no token in %s::%s (%v)", keyGordonServiceName, keyGordonGitHubToken, err)
		return &Config{}
	}

	return &Config{GitHub: GitHubConfig{Token: token}}
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
