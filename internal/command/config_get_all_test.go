package command_test

import (
	"github.com/kyoh86/gordon/internal/command"
	"github.com/kyoh86/gordon/internal/context"
)

func ExampleConfigGetAll() {
	if err := command.ConfigGetAll(&context.Config{
		GitHub: context.GitHubConfig{
			Token: "tokenx1",
			Host:  "hostx1",
			User:  "kyoh86",
		},
		Log: context.LogConfig{
			Level:        "trace",
			Date:         context.TrueOption,
			Time:         context.FalseOption,
			MicroSeconds: context.TrueOption,
			LongFile:     context.TrueOption,
			ShortFile:    context.TrueOption,
			UTC:          context.TrueOption,
		},
		History: context.HistoryConfig{
			File: "/var/log/gordon/history",
			Save: context.TrueOption,
		},
		Extract: context.ExtractConfig{
			Modes:   []context.FileMode{0111},
			Exclude: "exclude",
			Include: "include",
		},
		VRoot:         "/foo",
		VArchitecture: "arch",
		VOS:           "os",
	}); err != nil {
		panic(err)
	}
	// Unordered output:
	// log.level: trace
	// log.date: yes
	// log.time: no
	// log.microseconds: yes
	// log.longfile: yes
	// log.shortfile: yes
	// log.utc: yes
	// github.host: hostx1
	// github.user: kyoh86
	// github.token: *****
	// history.file: /var/log/gordon/history
	// history.save: yes
	// extract.modes: 0111
	// extract.exclude: exclude
	// extract.include: include
	// root: /foo
	// architecture: arch
	// os: os
}
