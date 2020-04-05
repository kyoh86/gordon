package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/kyoh86/gogh/gogh"
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

		download,

		get,
		install,
		// reinstall,
		// link,
		// relink,
		// unlink,
		uninstall,
		cleanup,

		dump,

		dump,
		restore,
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

func download(app *kingpin.Application) (string, func() error) {
	var (
		spec   gogh.RepoSpec
		update bool
		tag    string
	)
	cmd := app.Command("download", "Download from GitHub Release").Alias("download")
	cmd.Flag("update", "Update files").Short('u').BoolVar(&update)
	cmd.Flag("tag", "Target tag").StringVar(&tag)
	cmd.Arg("release", "Target repository (<repository URL> | <user>/<project> | <project>)").Required().SetValue(&spec)

	return mainutil.WrapCommand(cmd, func(ev command.Env) error {
		return command.Download(context.Background(), ev, spec, tag, update)
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
	cmd := app.Command("dump", "Dump installed versions")
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
