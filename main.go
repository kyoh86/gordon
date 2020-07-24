package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/kyoh86/gordon/internal/command"
	"github.com/kyoh86/gordon/internal/env"
	"github.com/kyoh86/gordon/internal/gordon"
	"github.com/kyoh86/gordon/internal/mainutil"
)

// nolint
var (
	version = "snapshot"
	commit  = "snapshot"
	date    = "snapshot"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	app := kingpin.New("gordon", "GO Released binaries DOwNloader").
		Version(fmt.Sprintf("%s-%s (%s)", version, commit, date)).Author("kyoh86")
	app.Command("config", "Get and set options")

	cmds := map[string]func() error{}
	for _, f := range []func(*kingpin.Application) (string, func() error){
		configGetAll,
		configGet,
		configSet,
		configUnset,
		setup,

		get,
		install,
		update,
		// reinstall,
		// link,
		// relink,
		// unlink,
		uninstall,
		cleanup,

		dump,
		restore,

		bin,
		initialize,
	} {
		key, run := f(app)
		cmds[key] = run
	}
	if err := cmds[kingpin.MustParse(app.Parse(os.Args[1:]))](); err != nil {
		log.Fatalf("error: %s", err)
	}
}

func configGetAll(app *kingpin.Application) (string, func() error) {
	cmd := app.GetCommand("config").Command("get-all", "get all options").Alias("list").Alias("ls")

	return mainutil.WrapConfigurableCommand(cmd, command.ConfigGetAll)
}

func configGet(app *kingpin.Application) (string, func() error) {
	var (
		name string
	)
	cmd := app.GetCommand("config").Command("get", "get an option")
	cmd.Arg("name", "option name").Required().StringVar(&name)

	return mainutil.WrapConfigurableCommand(cmd, func(_ command.Env, cfg *env.Config) error {
		return command.ConfigGet(cfg, name)
	})
}

func configSet(app *kingpin.Application) (string, func() error) {
	var (
		name  string
		value string
	)
	cmd := app.GetCommand("config").Command("set", "set an option")
	cmd.Arg("name", "option name").Required().StringVar(&name)
	cmd.Arg("value", "option value").Required().StringVar(&value)

	return mainutil.WrapConfigurableCommand(cmd, func(ev command.Env, cfg *env.Config) error {
		return command.ConfigSet(ev, cfg, name, value)
	})
}

func configUnset(app *kingpin.Application) (string, func() error) {
	var (
		name string
	)
	cmd := app.GetCommand("config").Command("unset", "unset an option").Alias("rm")
	cmd.Arg("name", "option name").Required().StringVar(&name)

	return mainutil.WrapConfigurableCommand(cmd, func(ev command.Env, cfg *env.Config) error {
		return command.ConfigUnset(ev, cfg, name)
	})
}

func setup(app *kingpin.Application) (string, func() error) {
	var (
		force bool
	)
	cmd := app.Command("setup", "Setup gordon with wizards")
	cmd.Flag("force", "Ask even though that the option has already set").BoolVar(&force)

	return mainutil.WrapConfigurableCommand(cmd, func(ev command.Env, cfg *env.Config) error {
		return command.Setup(context.Background(), ev, cfg, force)
	})
}

func get(app *kingpin.Application) (string, func() error) {
	var (
		spec gordon.VersionSpec
	)
	cmd := app.Command("get", "Download from GitHub Release")
	cmd.Arg("release", "Target release (<owner>/<name>[@<tag>])").Required().SetValue(&spec)

	return mainutil.WrapCommand(cmd, func(ev command.Env) error {
		return command.Get(context.Background(), ev, spec)
	})
}

func install(app *kingpin.Application) (string, func() error) {
	var (
		spec gordon.VersionSpec
	)
	cmd := app.Command("install", "Install from GitHub Release")
	cmd.Arg("release", "Target release (<owner>/<name>[@<tag>])").Required().SetValue(&spec)

	return mainutil.WrapCommand(cmd, func(ev command.Env) error {
		return command.Install(context.Background(), ev, spec)
	})
}

func uninstall(app *kingpin.Application) (string, func() error) {
	var (
		spec gordon.AppSpec
	)
	cmd := app.Command("uninstall", "Uninstall app")
	cmd.Arg("app", "Target app (<owner>/<name>)").Required().SetValue(&spec)

	return mainutil.WrapCommand(cmd, func(ev command.Env) error {
		return command.Uninstall(context.Background(), ev, spec)
	})
}

func dump(app *kingpin.Application) (string, func() error) {
	var (
		bundleFile string
	)
	cmd := app.Command("dump", "Dump installed versions").Alias("list")
	cmd.Arg("bundle-file", "Dumped version files").Default("-").StringVar(&bundleFile)
	return mainutil.WrapCommand(cmd, func(ev command.Env) error {
		return command.Dump(ev, bundleFile)
	})
}

func restore(app *kingpin.Application) (string, func() error) {
	var (
		bundleFile string
	)
	cmd := app.Command("restore", "Restore dumped versions")
	cmd.Arg("bundle-file", "Dumped version files").Default("-").StringVar(&bundleFile)
	return mainutil.WrapCommand(cmd, func(ev command.Env) error {
		return command.Restore(context.Background(), ev, bundleFile)
	})
}

func cleanup(app *kingpin.Application) (string, func() error) {
	cmd := app.Command("cleanup", "Clean cached versions")
	return mainutil.WrapCommand(cmd, func(ev command.Env) error {
		return command.Cleanup(ev)
	})
}

func update(app *kingpin.Application) (string, func() error) {
	cmd := app.Command("update", "Update installed applications")
	return mainutil.WrapCommand(cmd, func(ev command.Env) error {
		return command.Update(context.Background(), ev)
	})
}

func bin(app *kingpin.Application) (string, func() error) {
	cmd := app.Command("bin", "Print directory to store downloaded binaries")
	return mainutil.WrapCommand(cmd, func(ev command.Env) error {
		return command.Bin(context.Background(), ev)
	})
}

func initialize(app *kingpin.Application) (string, func() error) {
	cmd := app.Command("init", "Initialize shell to support gordon")
	return mainutil.WrapCommand(cmd, func(ev command.Env) error {
		return command.Initialize(context.Background(), ev)
	})
}
