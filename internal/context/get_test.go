package context

import (
	"bytes"
	"log"
	"os"
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
	assert.Equal(t, "", cfg.GitHubToken())
	assert.Equal(t, "github.com", cfg.GitHubHost())
	assert.Equal(t, "", cfg.GitHubUser())
	assert.Equal(t, "warn", cfg.LogLevel())
	assert.NotEmpty(t, cfg.Root())
	assert.NotEmpty(t, cfg.Arch())
	assert.NotEmpty(t, cfg.OS())
	assert.Equal(t, os.Stderr, cfg.Stderr())
	assert.Equal(t, os.Stdout, cfg.Stdout())
	assert.Equal(t, os.Stdin, cfg.Stdin())
}

func TestLoadConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		resetEnv(t)
		cfg, err := LoadConfig(bytes.NewBufferString(`
root: /foo
arch: arch
os: os

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
`))
		require.NoError(t, err)
		assert.Equal(t, "", cfg.GitHubToken(), "token should not be saved in file")
		assert.Equal(t, "hostx1", cfg.GitHubHost())
		assert.Equal(t, "kyoh86", cfg.GitHubUser())
		assert.Equal(t, "trace", cfg.LogLevel())
		assert.Equal(t, log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile|log.Lshortfile|log.LUTC, cfg.LogFlags())
		assert.Equal(t, "/foo", cfg.Root())
		assert.Equal(t, "arch", cfg.Arch())
		assert.Equal(t, "os", cfg.OS())
		assert.Equal(t, os.Stderr, cfg.Stderr())
		assert.Equal(t, os.Stdout, cfg.Stdout())
	})
	t.Run("invalid format", func(t *testing.T) {
		resetEnv(t)
		_, err := LoadConfig(bytes.NewBufferString(`{`))
		assert.Error(t, err)
	})
}

func TestSaveConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		resetEnv(t)

		var buf bytes.Buffer
		var cfg Config
		cfg.GitHub.Token = "token1"
		cfg.GitHub.Host = "hostx1"
		cfg.GitHub.User = "kyoh86"
		cfg.Log.Level = "trace"
		cfg.Log.Date = TrueOption
		cfg.Log.Time = FalseOption
		cfg.Log.MicroSeconds = TrueOption
		cfg.Log.LongFile = TrueOption
		cfg.Log.ShortFile = TrueOption
		cfg.Log.UTC = TrueOption
		cfg.VRoot = "/foo"
		cfg.VArch = "arch"
		cfg.VOS = "os"

		require.NoError(t, SaveConfig(&buf, &cfg))

		output := buf.String()
		assert.Contains(t, output, "root: /foo")
		assert.Contains(t, output, "arch: arch")
		assert.Contains(t, output, "os: os")
		assert.Contains(t, output, "log:")
		assert.Contains(t, output, "  level: trace")
		assert.Contains(t, output, `  date: "yes"`)
		assert.Contains(t, output, `  time: "no"`)
		assert.Contains(t, output, `  microseconds: "yes"`)
		assert.Contains(t, output, `  longfile: "yes"`)
		assert.Contains(t, output, `  shortfile: "yes"`)
		assert.Contains(t, output, `  utc: "yes"`)
		assert.Contains(t, output, "github:")
		assert.NotContains(t, output, "tokenx1")
		assert.Contains(t, output, "  user: kyoh86")
		assert.Contains(t, output, "  host: hostx1")
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
	require.NoError(t, os.Setenv(envGordonRoot, "/foo"))
	require.NoError(t, os.Setenv(envGordonArch, "arch"))
	require.NoError(t, os.Setenv(envGordonOS, "os"))
	cfg, err := GetEnvarConfig()
	require.NoError(t, err)
	assert.Equal(t, "tokenx1", cfg.GitHubToken())
	assert.Equal(t, "hostx1", cfg.GitHubHost())
	assert.Equal(t, "kyoh86", cfg.GitHubUser())
	assert.Equal(t, "trace", cfg.LogLevel())
	assert.Equal(t, log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile|log.Lshortfile|log.LUTC, cfg.LogFlags())
	assert.Equal(t, "/foo", cfg.Root(), "expects roots are not duplicated")
	assert.Equal(t, "arch", cfg.Arch())
	assert.Equal(t, "os", cfg.OS())
	assert.Equal(t, os.Stderr, cfg.Stderr())
	assert.Equal(t, os.Stdout, cfg.Stdout())
}
