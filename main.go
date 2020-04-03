package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/alecthomas/kingpin"
	"github.com/comail/colog"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gordon/internal/command"
	"github.com/kyoh86/gordon/internal/env"
	"github.com/kyoh86/xdg"
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
		configPut,
		configUnset,

		get,
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

	return wrapConfigurableCommand(cmd, command.ConfigGetAll)
}

func configGet(app *kingpin.Application) (string, func() error) {
	var (
		name string
	)
	cmd := app.GetCommand("config").Command("get", "get an option")
	cmd.Arg("name", "option name").Required().StringVar(&name)

	return wrapConfigurableCommand(cmd, func(_ env.Env, cfg *env.Config) error {
		return command.ConfigGet(cfg, name)
	})
}

func configPut(app *kingpin.Application) (string, func() error) {
	var (
		name  string
		value string
	)
	cmd := app.GetCommand("config").Command("put", "put an option").Alias("set")
	cmd.Arg("name", "option name").Required().StringVar(&name)
	cmd.Arg("value", "option value").Required().StringVar(&value)

	return wrapConfigurableCommand(cmd, func(ev env.Env, cfg *env.Config) error {
		return command.ConfigSet(ev, cfg, name, value)
	})
}

func configUnset(app *kingpin.Application) (string, func() error) {
	var (
		name string
	)
	cmd := app.GetCommand("config").Command("unset", "unset an option").Alias("rm")
	cmd.Arg("name", "option name").Required().StringVar(&name)

	return wrapConfigurableCommand(cmd, func(ev env.Env, cfg *env.Config) error {
		return command.ConfigUnset(ev, cfg, name)
	})
}

func get(app *kingpin.Application) (string, func() error) {
	var (
		spec   gogh.RepoSpec
		update bool
		tag    string
	)
	cmd := app.Command("get", "Clone/sync with a remote repository").Alias("download")
	cmd.Arg("repository", "Target repository (<repository URL> | <user>/<project> | <project>)").Required().SetValue(&spec)
	cmd.Flag("update", "Update files").Short('u').BoolVar(&update)
	cmd.Flag("tag", "Target tag").StringVar(&tag)

	return wrapCommand(cmd, func(ev env.Env) error {
		return command.Download(context.Background(), ev, spec, tag, update)
	})

}

func setConfigFlag(cmd *kingpin.CmdClause, configFile *string) {
	cmd.Flag("config", "configuration file").
		Default(filepath.Join(xdg.ConfigHome(), "gordon", "config.yaml")).
		Envar("GORDON_CONFIG").
		StringVar(configFile)
}

var plainLabels = colog.LevelMap{
	colog.LTrace:   []byte("[ trace ] "),
	colog.LDebug:   []byte("⚙ "),
	colog.LInfo:    []byte("ⓘ "),
	colog.LWarning: []byte("⚠ "),
	colog.LError:   []byte("☢ "),
	colog.LAlert:   []byte("☠ "),
}

var colorLabels = colog.LevelMap{
	colog.LTrace:   []byte("[ trace ] "),
	colog.LDebug:   []byte("\x1b[0;36m\u2699 \x1b[0m"),
	colog.LInfo:    []byte("\x1b[0;32m\u24d8 \x1b[0m"),
	colog.LWarning: []byte("\x1b[0;33m\u26a0 \x1b[0m"),
	colog.LError:   []byte("\x1b[0;31m\u2622 \x1b[0m"),
	colog.LAlert:   []byte("\x1b[0;37;41m\u2620 \x1b[0m"),
}

func openYAML(filename string) (io.Reader, func() error, error) {
	var reader io.Reader
	var teardown func() error
	file, err := os.Open(filename)
	switch {
	case err == nil:
		teardown = file.Close
		reader = file
	case os.IsNotExist(err):
		reader = env.EmptyYAMLReader
		teardown = func() error { return nil }
	default:
		return nil, nil, err
	}
	return reader, teardown, nil
}

func wrapCommand(cmd *kingpin.CmdClause, f func(env.Env) error) (string, func() error) {
	var configFile string
	setConfigFlag(cmd, &configFile)
	return cmd.FullCommand(), func() (retErr error) {
		reader, teardown, err := openYAML(configFile)
		if err != nil {
			return err
		}
		defer func() {
			if err := teardown(); err != nil && retErr == nil {
				retErr = err
				return
			}
		}()

		access, err := env.GetAccess(reader, env.EnvarPrefix)
		if err != nil {
			return err
		}

		return f(&access)
	}
}

func wrapConfigurableCommand(cmd *kingpin.CmdClause, f func(env.Env, *env.Config) error) (string, func() error) {
	var configFile string
	setConfigFlag(cmd, &configFile)
	return cmd.FullCommand(), func() (retErr error) {
		reader, teardown, err := openYAML(configFile)
		if err != nil {
			return err
		}
		defer func() {
			if err := teardown(); err != nil && retErr == nil {
				retErr = err
				return
			}
		}()

		config, access, err := env.GetAppenv(reader, env.EnvarPrefix)
		if err != nil {
			return err
		}

		if err = f(&access, &config); err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(configFile), 0744); err != nil {
			return err
		}
		file, err := os.OpenFile(configFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()
		return config.Save(file)
	}
}
