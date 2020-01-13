package context

import (
	"errors"

	"github.com/kyoh86/ask"
	"github.com/zalando/go-keyring"
)

var (
	ErrEmptyValue        = errors.New("empty value")
	ErrInvalidOptionName = errors.New("invalid option name")
)

type OptionAccessor struct {
	optionName string
	getter     func(cfg *Config) string
	putter     func(cfg *Config, value string) error
	unsetter   func(cfg *Config) error
	prompt     func(cfg *Config) error
}

func (a OptionAccessor) Get(cfg *Config) string              { return a.getter(cfg) }
func (a OptionAccessor) Put(cfg *Config, value string) error { return a.putter(cfg, value) }
func (a OptionAccessor) Unset(cfg *Config) error             { return a.unsetter(cfg) }
func (a OptionAccessor) Prompt(cfg *Config) error            { return a.prompt(cfg) }

var (
	configAccessor  map[string]OptionAccessor
	optionNames     []string
	optionAccessors = []OptionAccessor{
		logLevelOptionAccessor,
		logDateOptionAccessor,
		logTimeOptionAccessor,
		logMicroSecondsOptionAccessor,
		logLongFileOptionAccessor,
		logShortFileOptionAccessor,
		logUTCOptionAccessor,
		gitHubHostOptionAccessor,
		gitHubUserOptionAccessor,
		gitHubTokenOptionAccessor,
		historyFileOptionAccessor,
		historySaveOptionAccessor,
		extractModesOptionAccessor,
		extractExcludeOptionAccessor,
		extractIncludeOptionAccessor,
		rootOptionAccessor,
		architectureOptionAccessor,
		osOptionAccessor,
	}
)

func init() {
	m := map[string]OptionAccessor{}
	n := make([]string, 0, len(optionAccessors))
	for _, a := range optionAccessors {
		n = append(n, a.optionName)
		m[a.optionName] = a
	}
	configAccessor = m
	optionNames = n
}

func Option(optionName string) (*OptionAccessor, error) {
	a, ok := configAccessor[optionName]
	if !ok {
		return nil, ErrInvalidOptionName
	}
	return &a, nil
}

func OptionNames() []string {
	return optionNames
}

func Options() []OptionAccessor {
	return optionAccessors
}

var (
	logLevelOptionAccessor = OptionAccessor{
		optionName: "log.level",
		getter: func(cfg *Config) string {
			return cfg.LogLevel()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			if err := ValidateLogLevel(value); err != nil {
				return err
			}
			cfg.Log.Level = value
			return nil
		},
		unsetter: func(cfg *Config) error {
			cfg.Log.Level = ""
			return nil
		},
		prompt: func(cfg *Config) error {
			var value string
			doer := ask.Default(cfg.LogLevel()).Message("Logging severity levels").Enum([]string{
				"trace",
				"debug",
				"info",
				"warn",
				"error",
				"alert",
				"panic",
			}).Optional(true).Validation(ValidateLogLevel).StringVar(&value)
			switch err := doer.Do(); err {
			case nil:
				// noop
			case ask.ErrSkip:
				return nil
			default:
				return err
			}
			cfg.Log.Level = value
			return nil
		},
	}

	logDateOptionAccessor = OptionAccessor{
		optionName: "log.date",
		getter: func(cfg *Config) string {
			return cfg.Log.Date.String()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			return cfg.Log.Date.Decode(value)
		},
		unsetter: func(cfg *Config) error {
			cfg.Log.Date = EmptyBoolOption
			return nil
		},
		prompt: func(cfg *Config) error {
			return boolOptionPrompt(&cfg.Log.Date, "Logging the date in the local time zone like '2009/01/23'")
		},
	}

	logTimeOptionAccessor = OptionAccessor{
		optionName: "log.time",
		getter: func(cfg *Config) string {
			return cfg.Log.Time.String()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			return cfg.Log.Time.Decode(value)
		},
		unsetter: func(cfg *Config) error {
			cfg.Log.Time = EmptyBoolOption
			return nil
		},
		prompt: func(cfg *Config) error {
			return boolOptionPrompt(&cfg.Log.Time, "Logging the time in the local time zone like '01:23:23'")
		},
	}

	logMicroSecondsOptionAccessor = OptionAccessor{
		optionName: "log.microseconds",
		getter: func(cfg *Config) string {
			return cfg.Log.MicroSeconds.String()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			return cfg.Log.MicroSeconds.Decode(value)
		},
		unsetter: func(cfg *Config) error {
			cfg.Log.MicroSeconds = EmptyBoolOption
			return nil
		},
		prompt: func(cfg *Config) error {
			return boolOptionPrompt(&cfg.Log.MicroSeconds, "Logging microsecond resolution like '01:23:23.123123'.  Assumes Log.Time.")
		},
	}

	logLongFileOptionAccessor = OptionAccessor{
		optionName: "log.longfile",
		getter: func(cfg *Config) string {
			return cfg.Log.LongFile.String()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			return cfg.Log.LongFile.Decode(value)
		},
		unsetter: func(cfg *Config) error {
			cfg.Log.LongFile = EmptyBoolOption
			return nil
		},
		prompt: func(cfg *Config) error {
			return boolOptionPrompt(&cfg.Log.LongFile, "Logging full file name and line number: /a/b/c/d.go:23")
		},
	}

	logShortFileOptionAccessor = OptionAccessor{
		optionName: "log.shortfile",
		getter: func(cfg *Config) string {
			return cfg.Log.ShortFile.String()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			return cfg.Log.ShortFile.Decode(value)
		},
		unsetter: func(cfg *Config) error {
			cfg.Log.ShortFile = EmptyBoolOption
			return nil
		},
		prompt: func(cfg *Config) error {
			return boolOptionPrompt(&cfg.Log.ShortFile, "Logging final file name element and line number: d.go:23. overrides Log.Longfile")
		},
	}

	logUTCOptionAccessor = OptionAccessor{
		optionName: "log.utc",
		getter: func(cfg *Config) string {
			return cfg.Log.UTC.String()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			return cfg.Log.UTC.Decode(value)
		},
		unsetter: func(cfg *Config) error {
			cfg.Log.UTC = EmptyBoolOption
			return nil
		},
		prompt: func(cfg *Config) error {
			return boolOptionPrompt(&cfg.Log.UTC, "If Log.Date or Log.Time is set, use UTC rather than the local time zone")
		},
	}

	gitHubUserOptionAccessor = OptionAccessor{
		optionName: "github.user",
		getter: func(cfg *Config) string {
			return cfg.GitHubUser()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			if err := ValidateOwner(value); err != nil {
				return err
			}
			cfg.GitHub.User = value
			return nil
		},
		unsetter: func(cfg *Config) error {
			cfg.GitHub.User = ""
			return nil
		},
		prompt: func(cfg *Config) error {
			return stringOptionPrompt(cfg, "GitHub user name like 'kyoh86'", ValidateOwner, &cfg.GitHub.User)
		},
	}

	gitHubTokenOptionAccessor = OptionAccessor{
		optionName: "github.token",
		getter: func(cfg *Config) string {
			if cfg.GitHubToken() == "" {
				return ""
			}
			return "*****"
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			return keyring.Set(keyGordonServiceName, keyGordonGitHubToken, value)
		},
		unsetter: func(cfg *Config) error {
			return keyring.Delete(keyGordonServiceName, keyGordonGitHubToken)
		},
		prompt: func(cfg *Config) error {
			var value string
			doer := ask.Default(cfg.GitHubToken()).
				Message("GitHub API token").Optional(true).Hidden(true).StringVar(&value)
			switch err := doer.Do(); err {
			case nil:
				// noop
			case ask.ErrSkip:
				return nil
			default:
				return err
			}
			return keyring.Set(keyGordonServiceName, keyGordonGitHubToken, value)
		},
	}

	gitHubHostOptionAccessor = OptionAccessor{
		optionName: "github.host",
		getter: func(cfg *Config) string {
			return cfg.GitHubHost()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			cfg.GitHub.Host = value
			return nil
		},
		unsetter: func(cfg *Config) error {
			cfg.GitHub.Host = ""
			return nil
		},
		prompt: func(cfg *Config) error {
			return stringOptionPrompt(cfg, "GitHub host name like 'github.com'", nil, &cfg.GitHub.Host)
		},
	}

	historyFileOptionAccessor = OptionAccessor{
		optionName: "history.file",
		getter: func(cfg *Config) string {
			return cfg.History.File
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			if err := ValidateFile(value); err != nil {
				return err
			}
			cfg.History.File = value
			return nil
		},
		unsetter: func(cfg *Config) error {
			cfg.History.File = ""
			return nil
		},
		prompt: func(cfg *Config) error {
			return stringOptionPrompt(cfg, "A file to save downloading history", ValidateFile, &cfg.History.File)
		},
	}

	historySaveOptionAccessor = OptionAccessor{
		optionName: "history.save",
		getter: func(cfg *Config) string {
			return cfg.History.Save.String()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			return cfg.History.Save.Decode(value)
		},
		unsetter: func(cfg *Config) error {
			cfg.History.Save = EmptyBoolOption
			return nil
		},
		prompt: func(cfg *Config) error {
			return boolOptionPrompt(&cfg.History.Save, "Save downloading history to file")
		},
	}

	extractModesOptionAccessor = OptionAccessor{
		optionName: "extract.modes",
		getter: func(cfg *Config) string {
			buf, _ := cfg.Extract.Modes.MarshalText()
			return string(buf)
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			var m FileMode
			if err := m.UnmarshalText([]byte(value)); err != nil {
				return err
			}
			cfg.Extract.Modes = append(cfg.Extract.Modes, m)
			return nil
		},
		unsetter: func(cfg *Config) error {
			cfg.Extract.Modes = nil
			return nil
		},
		prompt: func(cfg *Config) error {
			text, err := cfg.Extract.Modes.MarshalText()
			if err != nil {
				return err
			}
			doer := ask.Default(string(text)).
				Message("File mode filters for files to extract from downloaded archive").
				Optional(true).Var(&cfg.Extract.Modes)
			switch err := doer.Do(); err {
			case ask.ErrSkip:
				return nil
			default:
				return err
			}
		},
	}

	extractExcludeOptionAccessor = OptionAccessor{
		optionName: "extract.exclude",
		getter: func(cfg *Config) string {
			return cfg.Extract.Exclude
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			if err := ValidateRegexp(value); err != nil {
				return err
			}
			cfg.Extract.Exclude = value
			return nil
		},
		unsetter: func(cfg *Config) error {
			cfg.Extract.Exclude = ""
			return nil
		},
		prompt: func(cfg *Config) error {
			return stringOptionPrompt(cfg, "File exclusion filter for files to extract from downloaded archive", ValidateRegexp, &cfg.Extract.Exclude)
		},
	}

	extractIncludeOptionAccessor = OptionAccessor{
		optionName: "extract.include",
		getter: func(cfg *Config) string {
			return cfg.Extract.Include
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			if err := ValidateRegexp(value); err != nil {
				return err
			}
			cfg.Extract.Include = value
			return nil
		},
		unsetter: func(cfg *Config) error {
			cfg.Extract.Include = ""
			return nil
		},
		prompt: func(cfg *Config) error {
			return stringOptionPrompt(cfg, "File inclusion filter for files to extract from downloaded archive", ValidateRegexp, &cfg.Extract.Exclude)
		},
	}

	rootOptionAccessor = OptionAccessor{
		optionName: "root",
		getter: func(cfg *Config) string {
			return cfg.Root()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			if err := ValidateRoot(value); err != nil {
				return err
			}
			cfg.VRoot = value
			return nil
		},
		unsetter: func(cfg *Config) error {
			cfg.VRoot = ""
			return nil
		},
		prompt: func(cfg *Config) error {
			return stringOptionPrompt(cfg, "Root directory to store downloaded files", ValidateRoot, &cfg.VRoot)
		},
	}

	architectureOptionAccessor = OptionAccessor{
		optionName: "architecture",
		getter: func(cfg *Config) string {
			return cfg.Architecture()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			if err := ValidateArchitecture(value); err != nil {
				return err
			}
			cfg.VArchitecture = value
			return nil
		},
		unsetter: func(cfg *Config) error {
			cfg.VArchitecture = ""
			return nil
		},
		prompt: func(cfg *Config) error {
			return stringOptionPrompt(cfg, "Target architecture", ValidateArchitecture, &cfg.VArchitecture)
		},
	}

	osOptionAccessor = OptionAccessor{
		optionName: "os",
		getter: func(cfg *Config) string {
			return cfg.OS()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			if err := ValidateOS(value); err != nil {
				return err
			}
			cfg.VOS = value
			return nil
		},
		unsetter: func(cfg *Config) error {
			cfg.VOS = ""
			return nil
		},
		prompt: func(cfg *Config) error {
			return stringOptionPrompt(cfg, "Target OS", ValidateOS, &cfg.VOS)
		},
	}
)

func boolOptionPrompt(opt *BoolOption, msg string) error {
	var value bool
	doer := ask.Default(opt.String()).Message(msg).Optional(true).YesNoVar(&value)
	switch err := doer.Do(); err {
	case nil:
		// noop
	case ask.ErrSkip:
		return nil
	default:
		return err
	}
	opt.SetBool(value)
	return nil
}

func stringOptionPrompt(msg string, validator func(string) error, target *string) error {
	var value string
	asker := ask.Default(*target).Message(msg).Optional(true)
	if validator != nil {
		asker = asker.Validation(validator)
	}
	switch err := asker.StringVar(&value).Do(); err {
	case nil:
		// noop
	case ask.ErrSkip:
		return nil
	default:
		return err
	}
	*target = value
	return nil
}
