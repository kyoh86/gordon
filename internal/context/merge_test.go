package context

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMergeConfig(t *testing.T) {
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

	cfg1, err := GetEnvarConfig()
	require.NoError(t, err)

	t.Run("full overwritten config", func(t *testing.T) {
		require.NoError(t, os.Setenv(envGordonGitHubToken, "tokenx2"))
		require.NoError(t, os.Setenv(envGordonGitHubHost, "hostx2"))
		require.NoError(t, os.Setenv(envGordonGitHubUser, "kyoh87"))
		require.NoError(t, os.Setenv(envGordonLogLevel, "debug"))
		require.NoError(t, os.Setenv(envGordonLogDate, "no"))
		require.NoError(t, os.Setenv(envGordonLogTime, "no"))
		require.NoError(t, os.Setenv(envGordonLogMicroSeconds, "no"))
		require.NoError(t, os.Setenv(envGordonLogLongFile, "no"))
		require.NoError(t, os.Setenv(envGordonLogShortFile, "no"))
		require.NoError(t, os.Setenv(envGordonLogUTC, "no"))
		require.NoError(t, os.Setenv(envGordonRoot, "/baz"))
		require.NoError(t, os.Setenv(envGordonArch, "arch"))
		require.NoError(t, os.Setenv(envGordonOS, "os"))

		cfg2, err := GetEnvarConfig()
		require.NoError(t, err)

		cfg := MergeConfig(cfg1, cfg2) // prior after config
		assert.Equal(t, "tokenx2", cfg.GitHubToken())
		assert.Equal(t, "hostx2", cfg.GitHubHost())
		assert.Equal(t, "kyoh87", cfg.GitHubUser())
		assert.Equal(t, "debug", cfg.LogLevel())
		assert.Equal(t, 0, cfg.LogFlags())
		assert.Equal(t, "/baz", cfg.Root())
		assert.Equal(t, "arch", cfg.Arch())
		assert.Equal(t, "os", cfg.OS())
		assert.Equal(t, os.Stderr, cfg.Stderr())
		assert.Equal(t, os.Stdout, cfg.Stdout())
	})

	t.Run("no overwritten config", func(t *testing.T) {
		resetEnv(t)

		cfg2, err := GetEnvarConfig()
		require.NoError(t, err)

		cfg := MergeConfig(cfg1, cfg2) // prior after config
		assert.Equal(t, "tokenx1", cfg.GitHubToken())
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

	resetEnv(t)
	assert.Equal(t, EmptyBoolOption, mergeBoolOption(EmptyBoolOption, EmptyBoolOption))
	assert.Equal(t, TrueOption, mergeBoolOption(TrueOption, EmptyBoolOption))
	assert.Equal(t, FalseOption, mergeBoolOption(FalseOption, EmptyBoolOption))
	assert.Equal(t, TrueOption, mergeBoolOption(EmptyBoolOption, TrueOption))
	assert.Equal(t, FalseOption, mergeBoolOption(TrueOption, FalseOption))
	assert.Equal(t, TrueOption, mergeBoolOption(FalseOption, TrueOption))
}
