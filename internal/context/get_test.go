package context

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetEnv(t *testing.T) {
	t.Helper()
	for _, key := range envNames {
		require.NoError(t, os.Setenv(key, ""))
	}
}

func TestDefaultConfig(t *testing.T) {
	resetEnv(t)
	cfg := DefaultConfig()
	assert.Equal(t, os.Stderr, cfg.Stderr())
	assert.Equal(t, os.Stdout, cfg.Stdout())
	assert.Equal(t, os.Stdin, cfg.Stdin())
	assert.Equal(t, "warn", cfg.LogLevel())
	assert.Equal(t, log.Ltime, cfg.LogFlags())
	assert.False(t, cfg.LogDate())
	assert.True(t, cfg.LogTime())
	assert.False(t, cfg.LogMicroSeconds())
	assert.False(t, cfg.LogLongFile())
	assert.False(t, cfg.LogShortFile())
	assert.False(t, cfg.LogUTC())
	assert.Equal(t, "", cfg.GitHubToken())
	assert.Equal(t, "github.com", cfg.GitHubHost())
	assert.Equal(t, "", cfg.GitHubUser())
	assert.NotEmpty(t, cfg.HistoryFile())
	assert.True(t, strings.HasSuffix(
		cfg.HistoryFile(),
		strings.Join([]string{"", "gordon", "history"}, string(filepath.Separator)),
	), cfg.HistoryFile())
	assert.True(t, cfg.HistorySave())
	assert.EqualValues(t, FileModes{0111}, cfg.ExtractModes())
	assert.Empty(t, cfg.ExtractExclude())
	assert.Empty(t, cfg.ExtractInclude())
	assert.NotEmpty(t, cfg.Root())
	assert.NotEmpty(t, cfg.Architecture())
	assert.NotEmpty(t, cfg.OS())
}

func TestLoadConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		resetEnv(t)
		cfg, err := LoadConfig(bytes.NewBufferString(`
log:
  level: trace
  date: "yes"
  time: "yes"
  microseconds: yes
  longfile: "yes"
  shortfile: "yes"
  utc: "yes"

github:
  token: tokenx1
  user: kyoh86
  host: hostx1

history:
  file: history-file
  save: yes

extract:
  modes: 0755
  exclude: exclude
  include: include

root: /foo
architecture: arch
os: os
`))
		require.NoError(t, err)
		assert.Equal(t, "", cfg.GitHubToken(), "token should not be saved in file")
		assert.Equal(t, os.Stderr, cfg.Stderr())
		assert.Equal(t, os.Stdout, cfg.Stdout())
		assert.Equal(t, "trace", cfg.LogLevel())
		assert.Equal(t, log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile|log.Lshortfile|log.LUTC, cfg.LogFlags())
		assert.True(t, cfg.LogDate())
		assert.True(t, cfg.LogTime())
		assert.True(t, cfg.LogMicroSeconds())
		assert.True(t, cfg.LogLongFile())
		assert.True(t, cfg.LogShortFile())
		assert.True(t, cfg.LogUTC())
		assert.Equal(t, "hostx1", cfg.GitHubHost())
		assert.Equal(t, "kyoh86", cfg.GitHubUser())
		assert.Equal(t, "history-file", cfg.HistoryFile())
		assert.True(t, cfg.HistorySave())
		assert.EqualValues(t, FileModes{0755}, cfg.ExtractModes())
		assert.Equal(t, "exclude", cfg.ExtractExclude())
		assert.Equal(t, "include", cfg.ExtractInclude())
		assert.Equal(t, "/foo", cfg.Root())
		assert.Equal(t, "arch", cfg.Architecture())
		assert.Equal(t, "os", cfg.OS())
	})
}

func TestSaveConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		resetEnv(t)

		var buf bytes.Buffer
		var cfg Config
		cfg.Log.Level = "trace"
		cfg.Log.Date = TrueOption
		cfg.Log.Time = FalseOption
		cfg.Log.MicroSeconds = TrueOption
		cfg.Log.LongFile = TrueOption
		cfg.Log.ShortFile = TrueOption
		cfg.Log.UTC = TrueOption
		cfg.GitHub.Token = "token1"
		cfg.GitHub.Host = "hostx1"
		cfg.GitHub.User = "kyoh86"
		cfg.History.File = "history-file"
		cfg.History.Save = TrueOption
		cfg.Extract.Modes = FileModes{0111, 0222}
		cfg.Extract.Exclude = "exclude"
		cfg.Extract.Include = "include"
		cfg.VRoot = "/foo"
		cfg.VArchitecture = "arch"
		cfg.VOS = "os"

		require.NoError(t, SaveConfig(&buf, &cfg))

		output := buf.String()
		assert.Contains(t, output, "log:")
		assert.Contains(t, output, "  level: trace")
		assert.Contains(t, output, `  date: yes`)
		assert.Contains(t, output, `  time: no`)
		assert.Contains(t, output, `  microseconds: yes`)
		assert.Contains(t, output, `  longfile: yes`)
		assert.Contains(t, output, `  shortfile: yes`)
		assert.Contains(t, output, `  utc: yes`)
		assert.Contains(t, output, "github:")
		assert.NotContains(t, output, "tokenx1")
		assert.Contains(t, output, "  user: kyoh86")
		assert.Contains(t, output, "  host: hostx1")
		assert.Contains(t, output, "history:")
		assert.Contains(t, output, "  file: history-file")
		assert.Contains(t, output, "  save: yes")
		assert.Contains(t, output, "extract:")
		assert.Contains(t, output, "  modes: 0111|0222")
		assert.Contains(t, output, "  exclude: exclude")
		assert.Contains(t, output, "  include: include")
		assert.Contains(t, output, "root: /foo")
		assert.Contains(t, output, "architecture: arch")
		assert.Contains(t, output, "os: os")
	})
}

func TestGetEnvarConfig(t *testing.T) {
	resetEnv(t)
	require.NoError(t, os.Setenv(envGordonGitHubToken, "tokenx1"))
	require.NoError(t, os.Setenv(envGordonGitHubHost, "hostx1"))
	require.NoError(t, os.Setenv(envGordonGitHubUser, "kyoh86"))
	require.NoError(t, os.Setenv(envGordonLogLevel, "trace"))
	require.NoError(t, os.Setenv(envGordonLogDate, "yes"))
	require.NoError(t, os.Setenv(envGordonLogTime, "yes"))
	require.NoError(t, os.Setenv(envGordonLogMicroSeconds, "yes"))
	require.NoError(t, os.Setenv(envGordonLogLongFile, "yes"))
	require.NoError(t, os.Setenv(envGordonLogShortFile, "yes"))
	require.NoError(t, os.Setenv(envGordonLogUTC, "yes"))
	require.NoError(t, os.Setenv(envGordonHistoryFile, "history-file"))
	require.NoError(t, os.Setenv(envGordonHistorySave, "yes"))
	require.NoError(t, os.Setenv(envGordonExtractModes, "0111|0222"))
	require.NoError(t, os.Setenv(envGordonExtractExclude, "exclude"))
	require.NoError(t, os.Setenv(envGordonExtractInclude, "include"))
	require.NoError(t, os.Setenv(envGordonRoot, "/foo"))
	require.NoError(t, os.Setenv(envGordonArchitecture, "arch"))
	require.NoError(t, os.Setenv(envGordonOS, "os"))
	cfg, err := GetEnvarConfig()
	require.NoError(t, err)
	assert.Equal(t, os.Stderr, cfg.Stderr())
	assert.Equal(t, os.Stdout, cfg.Stdout())
	assert.Equal(t, "tokenx1", cfg.GitHubToken())
	assert.Equal(t, "hostx1", cfg.GitHubHost())
	assert.Equal(t, "kyoh86", cfg.GitHubUser())
	assert.Equal(t, "trace", cfg.LogLevel())
	assert.Equal(t, log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile|log.Lshortfile|log.LUTC, cfg.LogFlags())
	assert.True(t, cfg.LogDate())
	assert.True(t, cfg.LogTime())
	assert.True(t, cfg.LogMicroSeconds())
	assert.True(t, cfg.LogLongFile())
	assert.True(t, cfg.LogShortFile())
	assert.True(t, cfg.LogUTC())
	assert.Equal(t, "history-file", cfg.HistoryFile())
	assert.True(t, cfg.HistorySave())
	assert.Equal(t, FileModes{0111, 0222}, cfg.ExtractModes())
	assert.Equal(t, "exclude", cfg.ExtractExclude())
	assert.Equal(t, "include", cfg.ExtractInclude())
	assert.Equal(t, "/foo", cfg.Root(), "expects roots are not duplicated")
	assert.Equal(t, "arch", cfg.Architecture())
	assert.Equal(t, "os", cfg.OS())
}
