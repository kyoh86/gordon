package context

import (
	"context"
	"io"
	"log"
	"os"
)

// Config holds configuration file values.
type Config struct {
	context.Context `yaml:"-"`
	Log             LogConfig     `yaml:"log,omitempty"`
	GitHub          GitHubConfig  `yaml:"github,omitempty"`
	History         HistoryConfig `yaml:"history,omitempty"`
	Extract         ExtractConfig `yaml:"extract,omitempty"`
	VRoot           string        `yaml:"root,omitempty" env:"GORDON_ROOT"`
	VArchitecture   string        `yaml:"architecture,omitempty" env:"GORDON_ARCHITECTURE"`
	VOS             string        `yaml:"os,omitempty" env:"GORDON_OS"`
}

type LogConfig struct {
	Level        string     `yaml:"level,omitempty" env:"GORDON_LOG_LEVEL"`
	Date         BoolOption `yaml:"date" env:"GORDON_LOG_DATE"`                 // the date in the local time zone: 2009/01/23
	Time         BoolOption `yaml:"time" env:"GORDON_LOG_TIME"`                 // the time in the local time zone: 01:23:23
	MicroSeconds BoolOption `yaml:"microseconds" env:"GORDON_LOG_MICROSECONDS"` // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	LongFile     BoolOption `yaml:"longfile" env:"GORDON_LOG_LONGFILE"`         // full file name and line number: /a/b/c/d.go:23
	ShortFile    BoolOption `yaml:"shortfile" env:"GORDON_LOG_SHORTFILE"`       // final file name element and line number: d.go:23. overrides Llongfile
	UTC          BoolOption `yaml:"utc" env:"GORDON_LOG_UTC"`                   // if Ldate or Ltime is set, use UTC rather than the local time zone
}

type GitHubConfig struct {
	Token string `yaml:"-" env:"GORDON_GITHUB_TOKEN"`
	User  string `yaml:"user,omitempty" env:"GORDON_GITHUB_USER"`
	Host  string `yaml:"host,omitempty" env:"GORDON_GITHUB_HOST"`
}

type HistoryConfig struct {
	File string     `yaml:"file,omitempty" env:"GORDON_HISTORY_FILE"`
	Save BoolOption `yaml:"save" env:"GORDON_HISTORY_SAVE"`
}

type ExtractConfig struct {
	Modes   FileModes `yaml:"modes,omitempty" env:"GORDON_EXTRACT_MODES"`
	Exclude string     `yaml:"exclude,omitempty" env:"GORDON_EXTRACT_EXCLUDE"`
	Include string     `yaml:"include,omitempty" env:"GORDON_EXTRACT_INCLUDE"`
}

func (c *Config) Stdin() io.Reader  { return os.Stdin }
func (c *Config) Stdout() io.Writer { return os.Stdout }
func (c *Config) Stderr() io.Writer { return os.Stderr }

func (c *Config) LogLevel() string { return c.Log.Level }

func (c *Config) LogFlags() int {
	var f int
	if c.Log.Date.Bool() {
		f |= log.Ldate
	}
	if c.Log.Time.Bool() {
		f |= log.Ltime
	}
	if c.Log.MicroSeconds.Bool() {
		f |= log.Lmicroseconds
	}
	if c.Log.LongFile.Bool() {
		f |= log.Llongfile
	}
	if c.Log.ShortFile.Bool() {
		f |= log.Lshortfile
	}
	if c.Log.UTC.Bool() {
		f |= log.LUTC
	}
	return f
}

func (c *Config) LogDate() bool         { return c.Log.Date.Bool() }
func (c *Config) LogTime() bool         { return c.Log.Time.Bool() }
func (c *Config) LogMicroSeconds() bool { return c.Log.MicroSeconds.Bool() }
func (c *Config) LogLongFile() bool     { return c.Log.LongFile.Bool() }
func (c *Config) LogShortFile() bool    { return c.Log.ShortFile.Bool() }
func (c *Config) LogUTC() bool          { return c.Log.UTC.Bool() }

func (c *Config) GitHubUser() string  { return c.GitHub.User }
func (c *Config) GitHubToken() string { return c.GitHub.Token }
func (c *Config) GitHubHost() string  { return c.GitHub.Host }

func (c *Config) HistoryFile() string { return c.History.File }
func (c *Config) HistorySave() bool   { return c.History.Save.Bool() }

func (c *Config) ExtractModes() FileModes { return c.Extract.Modes }
func (c *Config) ExtractExclude() string   { return c.Extract.Exclude }
func (c *Config) ExtractInclude() string   { return c.Extract.Include }

func (c *Config) Root() string         { return expandPath(c.VRoot) }
func (c *Config) Architecture() string { return c.VArchitecture }
func (c *Config) OS() string           { return c.VOS }
