package context

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
)

type MockContext struct {
	context.Context
	MStdin           io.Reader
	MStdout          io.Writer
	MStderr          io.Writer
	MLogLevel        string
	MLogFlags        int
	MLogDate         bool
	MLogTime         bool
	MLogMicroSeconds bool
	MLogLongFile     bool
	MLogShortFile    bool
	MLogUTC          bool
	MGitHubUser      string
	MGitHubToken     string
	MGitHubHost      string
	MHistoryFile     string
	MHistorySave     bool
	MExtractModes    FileModes
	MExtractExclude  string
	MExtractInclude  string
	MRoot            string
	MArchitecture    string
	MOS              string
}

func (c *MockContext) Stdin() io.Reader {
	if r := c.MStdin; r != nil {
		return r
	}
	return &bytes.Buffer{}
}

func (c *MockContext) Stdout() io.Writer {
	if w := c.MStdout; w != nil {
		return w
	}
	return ioutil.Discard
}

func (c *MockContext) Stderr() io.Writer {
	if w := c.MStderr; w != nil {
		return w
	}
	return ioutil.Discard
}

func (c *MockContext) LogLevel() string        { return c.MLogLevel }
func (c *MockContext) LogFlags() int           { return c.MLogFlags }
func (c *MockContext) LogDate() bool           { return c.MLogDate }
func (c *MockContext) LogTime() bool           { return c.MLogTime }
func (c *MockContext) LogMicroSeconds() bool   { return c.MLogMicroSeconds }
func (c *MockContext) LogLongFile() bool       { return c.MLogLongFile }
func (c *MockContext) LogShortFile() bool      { return c.MLogShortFile }
func (c *MockContext) LogUTC() bool            { return c.MLogUTC }
func (c *MockContext) GitHubUser() string      { return c.MGitHubUser }
func (c *MockContext) GitHubToken() string     { return c.MGitHubToken }
func (c *MockContext) GitHubHost() string      { return c.MGitHubHost }
func (c *MockContext) HistoryFile() string     { return c.MHistoryFile }
func (c *MockContext) HistorySave() bool       { return c.MHistorySave }
func (c *MockContext) ExtractModes() FileModes { return c.MExtractModes }
func (c *MockContext) ExtractExclude() string  { return c.MExtractExclude }
func (c *MockContext) ExtractInclude() string  { return c.MExtractInclude }
func (c *MockContext) Root() string            { return c.MRoot }
func (c *MockContext) Architecture() string    { return c.MArchitecture }
func (c *MockContext) OS() string              { return c.MOS }
