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
		VRoot: "/foo",
		VArch: "arch",
		VOS:   "os",
	}); err != nil {
		panic(err)
	}
	// Unordered output:
	// root = /foo
	// arch = arch
	// os = os
	// github.host = hostx1
	// github.user = kyoh86
	// github.token = *****
	// log.level = trace
	// log.date = yes
	// log.time = no
	// log.microseconds = yes
	// log.longfile = yes
	// log.shortfile = yes
	// log.utc = yes
}
