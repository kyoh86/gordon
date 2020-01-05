package context

import (
	"context"
	"io"
)

// Context holds configurations and environments
type Context interface {
	context.Context
	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
	GitHubUser() string
	GitHubToken() string
	GitHubHost() string
	Root() string
	Arch() string
	OS() string
	LogLevel() string
	LogFlags() int // log.Lxxx flags
	LogDate() bool
	LogTime() bool
	LogMicroSeconds() bool
	LogLongFile() bool
	LogShortFile() bool
	LogUTC() bool
}
