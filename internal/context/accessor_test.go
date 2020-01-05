package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccessor(t *testing.T) {
	const (
		dummyToken = "token1"
		dummyHost  = "hostx1"
		dummyUser  = "kyoh86"
		dummyLevel = "trace"
	)
	t.Run("getting", func(t *testing.T) {
		mustOption := func(acc *OptionAccessor, err error) *OptionAccessor {
			t.Helper()
			require.NoError(t, err)
			return acc
		}
		var cfg Config
		cfg.Log.Level = dummyLevel
		cfg.Log.Date = TrueOption
		cfg.Log.Time = FalseOption
		cfg.Log.MicroSeconds = TrueOption
		cfg.Log.LongFile = TrueOption
		cfg.Log.ShortFile = TrueOption
		cfg.Log.UTC = TrueOption
		cfg.GitHub.Token = dummyToken
		cfg.GitHub.Host = dummyHost
		cfg.GitHub.User = dummyUser
		cfg.History.File = "history-file"
		cfg.History.Save = TrueOption
		cfg.Extract.Modes = FileModes{0111, 0222}
		cfg.Extract.Exclude = "exclude"
		cfg.Extract.Include = "include"
		cfg.VRoot = "/foo"
		cfg.VArchitecture = "386"
		cfg.VOS = "linux"

		_, err := Option("invalid name")
		assert.EqualError(t, err, "invalid option name")
		assert.Equal(t, dummyLevel, mustOption(Option("log.level")).Get(&cfg))
		assert.Equal(t, "yes", mustOption(Option("log.date")).Get(&cfg))
		assert.Equal(t, "no", mustOption(Option("log.time")).Get(&cfg))
		assert.Equal(t, "yes", mustOption(Option("log.microseconds")).Get(&cfg))
		assert.Equal(t, "yes", mustOption(Option("log.longfile")).Get(&cfg))
		assert.Equal(t, "yes", mustOption(Option("log.shortfile")).Get(&cfg))
		assert.Equal(t, "yes", mustOption(Option("log.utc")).Get(&cfg))
		assert.Equal(t, "*****", mustOption(Option("github.token")).Get(&cfg))
		assert.Equal(t, dummyHost, mustOption(Option("github.host")).Get(&cfg))
		assert.Equal(t, dummyUser, mustOption(Option("github.user")).Get(&cfg))
		assert.Equal(t, "history-file", mustOption(Option("history.file")).Get(&cfg))
		assert.Equal(t, "yes", mustOption(Option("history.save")).Get(&cfg))
		assert.Equal(t, "0111|0222", mustOption(Option("extract.modes")).Get(&cfg))
		assert.Equal(t, "exclude", mustOption(Option("extract.exclude")).Get(&cfg))
		assert.Equal(t, "include", mustOption(Option("extract.include")).Get(&cfg))
		assert.Equal(t, "/foo", mustOption(Option("root")).Get(&cfg))
		assert.Equal(t, "386", mustOption(Option("architecture")).Get(&cfg))
		assert.Equal(t, "linux", mustOption(Option("os")).Get(&cfg))
	})

	t.Run("putting", func(t *testing.T) {
		mustOption := func(acc *OptionAccessor, err error) *OptionAccessor {
			t.Helper()
			require.NoError(t, err)
			return acc
		}
		var cfg Config
		assert.NoError(t, mustOption(Option("log.level")).Put(&cfg, dummyLevel))
		assert.NoError(t, mustOption(Option("log.date")).Put(&cfg, "yes"))
		assert.NoError(t, mustOption(Option("log.time")).Put(&cfg, "no"))
		assert.NoError(t, mustOption(Option("log.microseconds")).Put(&cfg, "yes"))
		assert.NoError(t, mustOption(Option("log.longfile")).Put(&cfg, "yes"))
		assert.NoError(t, mustOption(Option("log.shortfile")).Put(&cfg, "yes"))
		assert.NoError(t, mustOption(Option("log.utc")).Put(&cfg, "yes"))
		//TODO: put github.token test (change key name and service name for test)
		assert.NoError(t, mustOption(Option("github.host")).Put(&cfg, dummyHost))
		assert.NoError(t, mustOption(Option("github.user")).Put(&cfg, dummyUser))
		assert.NoError(t, mustOption(Option("history.file")).Put(&cfg, "history-file"))
		assert.NoError(t, mustOption(Option("history.save")).Put(&cfg, "yes"))
		assert.NoError(t, mustOption(Option("extract.modes")).Put(&cfg, "0111"))
		assert.NoError(t, mustOption(Option("extract.modes")).Put(&cfg, "0222"))
		assert.NoError(t, mustOption(Option("extract.exclude")).Put(&cfg, "exclude"))
		assert.NoError(t, mustOption(Option("extract.include")).Put(&cfg, "include"))
		assert.NoError(t, mustOption(Option("root")).Put(&cfg, "/foo"))
		assert.NoError(t, mustOption(Option("architecture")).Put(&cfg, "386"))
		assert.NoError(t, mustOption(Option("os")).Put(&cfg, "linux"))

		assert.Equal(t, dummyLevel, cfg.Log.Level)
		assert.Equal(t, TrueOption, cfg.Log.Date)
		assert.True(t, cfg.LogDate())
		assert.Equal(t, FalseOption, cfg.Log.Time)
		assert.False(t, cfg.LogTime())
		assert.Equal(t, TrueOption, cfg.Log.MicroSeconds)
		assert.True(t, cfg.LogMicroSeconds())
		assert.Equal(t, TrueOption, cfg.Log.LongFile)
		assert.True(t, cfg.LogLongFile())
		assert.Equal(t, TrueOption, cfg.Log.ShortFile)
		assert.True(t, cfg.LogShortFile())
		assert.Equal(t, TrueOption, cfg.Log.UTC)
		assert.True(t, cfg.LogUTC())
		assert.Equal(t, "", cfg.GitHub.Token)
		assert.Equal(t, dummyHost, cfg.GitHub.Host)
		assert.Equal(t, dummyUser, cfg.GitHub.User)
		assert.Equal(t, "history-file", cfg.History.File)
		assert.Equal(t, TrueOption, cfg.History.Save)
		assert.EqualValues(t, FileModes{0111, 0222}, cfg.Extract.Modes)
		assert.Equal(t, "exclude", cfg.Extract.Exclude)
		assert.Equal(t, "include", cfg.Extract.Include)
		assert.Equal(t, "/foo", cfg.VRoot)
		assert.Equal(t, "386", cfg.VArchitecture)
		assert.Equal(t, "linux", cfg.VOS)
	})
	t.Run("putting error", func(t *testing.T) {
		mustOption := func(acc *OptionAccessor, err error) *OptionAccessor {
			t.Helper()
			require.NoError(t, err)
			return acc
		}
		var cfg Config
		assert.EqualError(t, mustOption(Option("log.level")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("log.date")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("log.time")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("log.microseconds")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("log.longfile")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("log.shortfile")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("log.utc")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("github.host")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("github.user")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("history.file")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("history.save")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("extract.modes")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("extract.exclude")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("extract.include")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("root")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("architecture")).Put(&cfg, ""), "empty value")
		assert.EqualError(t, mustOption(Option("os")).Put(&cfg, ""), "empty value")

		assert.Error(t, mustOption(Option("log.level")).Put(&cfg, "foobar"), "invalid log level")
		assert.Error(t, mustOption(Option("log.date")).Put(&cfg, "invalid value"), "invalid value")
		assert.Error(t, mustOption(Option("log.time")).Put(&cfg, "invalid value"), "invalid value")
		assert.Error(t, mustOption(Option("log.microseconds")).Put(&cfg, "invalid value"), "invalid value")
		assert.Error(t, mustOption(Option("log.longfile")).Put(&cfg, "invalid value"), "invalid value")
		assert.Error(t, mustOption(Option("log.shortfile")).Put(&cfg, "invalid value"), "invalid value")
		assert.Error(t, mustOption(Option("log.utc")).Put(&cfg, "invalid value"), "invalid value")
		assert.Error(t, mustOption(Option("github.user")).Put(&cfg, "-kyoh86"), "invalid github username")
		assert.Error(t, mustOption(Option("history.save")).Put(&cfg, "invalid value"), "invalid value")
		assert.Error(t, mustOption(Option("extract.modes")).Put(&cfg, "x"), "invalid value")

		assert.Equal(t, "", cfg.Log.Level)
		assert.Equal(t, EmptyBoolOption, cfg.Log.Date)
		assert.Equal(t, EmptyBoolOption, cfg.Log.Time)
		assert.Equal(t, EmptyBoolOption, cfg.Log.MicroSeconds)
		assert.Equal(t, EmptyBoolOption, cfg.Log.LongFile)
		assert.Equal(t, EmptyBoolOption, cfg.Log.ShortFile)
		assert.Equal(t, EmptyBoolOption, cfg.Log.UTC)
		assert.Equal(t, "", cfg.GitHub.Token)
		assert.Equal(t, "", cfg.GitHub.Host)
		assert.Equal(t, "", cfg.GitHub.User)
		assert.Empty(t, cfg.History.File)
		assert.Equal(t, EmptyBoolOption, cfg.History.Save)
		assert.Empty(t, cfg.Extract.Modes)
		assert.Empty(t, cfg.Extract.Exclude)
		assert.Empty(t, cfg.Extract.Include)
		assert.Empty(t, cfg.VRoot)
		assert.Empty(t, cfg.VArchitecture)
		assert.Empty(t, cfg.VOS)
	})

	t.Run("unsetting", func(t *testing.T) {
		var cfg Config
		cfg.Log.Level = dummyLevel
		cfg.Log.Date = TrueOption
		cfg.Log.Time = FalseOption
		cfg.Log.MicroSeconds = TrueOption
		cfg.Log.LongFile = TrueOption
		cfg.Log.ShortFile = TrueOption
		cfg.Log.UTC = TrueOption
		// TODO: unset github.token test (change key name and service name for test)
		cfg.GitHub.Host = dummyHost
		cfg.GitHub.User = dummyUser
		cfg.History.File = "history-file"
		cfg.History.Save = TrueOption
		cfg.Extract.Modes = FileModes{0111}
		cfg.Extract.Exclude = "exclude"
		cfg.Extract.Include = "include"
		cfg.VRoot = "/foo"
		cfg.VArchitecture = "arch"
		cfg.VOS = "os"

		_, err := Option("invalid name")
		assert.EqualError(t, err, "invalid option name")
		for _, name := range OptionNames() {
			if name == gitHubTokenOptionAccessor.optionName {
				continue
			}
			acc, err := Option(name)
			require.NoError(t, err)
			assert.NoError(t, acc.Unset(&cfg), name)
		}
		assert.Equal(t, "", cfg.Log.Level)
		assert.Equal(t, EmptyBoolOption, cfg.Log.Date)
		assert.Equal(t, EmptyBoolOption, cfg.Log.Time)
		assert.Equal(t, EmptyBoolOption, cfg.Log.MicroSeconds)
		assert.Equal(t, EmptyBoolOption, cfg.Log.LongFile)
		assert.Equal(t, EmptyBoolOption, cfg.Log.ShortFile)
		assert.Equal(t, EmptyBoolOption, cfg.Log.UTC)
		assert.Equal(t, "", cfg.GitHub.Token)
		assert.Equal(t, "", cfg.GitHub.Host)
		assert.Equal(t, "", cfg.GitHub.User)
		assert.Empty(t, cfg.History.File)
		assert.Equal(t, EmptyBoolOption, cfg.History.Save)
		assert.Empty(t, cfg.Extract.Modes)
		assert.Empty(t, cfg.Extract.Exclude)
		assert.Empty(t, cfg.Extract.Include)
		assert.Empty(t, cfg.VRoot)
		assert.Empty(t, cfg.VArchitecture)
		assert.Empty(t, cfg.VOS)
	})
}
