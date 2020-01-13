package context

import (
	"context"
	"io"
)

// Context holds configurations and environments
type Context interface {
	context.Context
	IOContext
	LogContext
	GitHubContext
	HistoryContext
	ExtractContext
	Root() string
	Architecture() string
	OS() string
}
type IOContext interface {
	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
}
type LogContext interface {
	LogLevel() string
	LogFlags() int // log.Lxxx flags
	LogDate() bool
	LogTime() bool
	LogMicroSeconds() bool
	LogLongFile() bool
	LogShortFile() bool
	LogUTC() bool
}
type GitHubContext interface {
	GitHubUser() string
	GitHubToken() string
	GitHubHost() string
}
type HistoryContext interface {
	HistoryFile() string
	HistorySave() bool
}
type ExtractContext interface {
	ExtractModes() FileModes
	ExtractExclude() string
	ExtractInclude() string
}
