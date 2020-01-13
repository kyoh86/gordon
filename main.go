package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/alecthomas/kingpin"
	"github.com/comail/colog"
	"github.com/kyoh86/gogh/gogh"
	"github.com/kyoh86/gordon/internal/command"
	"github.com/kyoh86/gordon/internal/context"
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
		configSetup,

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

	return wrapConfigurableCommand(cmd, func(cfg *context.Config) error {
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

	return wrapConfigurableCommand(cmd, func(cfg *context.Config) error {
		return command.ConfigPut(cfg, name, value)
	})
}

func configUnset(app *kingpin.Application) (string, func() error) {
	var (
		name string
	)
	cmd := app.GetCommand("config").Command("unset", "unset an option").Alias("rm")
	cmd.Arg("name", "option name").Required().StringVar(&name)

	return wrapConfigurableCommand(cmd, func(cfg *context.Config) error {
		return command.ConfigUnset(cfg, name)
	})
}

func configSetup(app *kingpin.Application) (string, func() error) {
	cmd := app.GetCommand("config").Command("setup", "setup all options").Alias("prompt")

	return wrapConfigurableCommand(cmd, func(cfg *context.Config) error {
		return command.ConfigPrompt(cfg)
	})
}

func get(app *kingpin.Application) (string, func() error) {
	var (
		repo gogh.Repo
		tag  string
	)
	cmd := app.Command("get", "Clone/sync with a remote repository").Alias("download")
	cmd.Arg("repository", "Target repository (<repository URL> | <user>/<project> | <project>)").Required().SetValue(&repo)
	cmd.Flag("tag", "Target tag").StringVar(&tag)

	return wrapCommand(cmd, func(ctx context.Context) error {
		return command.Download(ctx, &repo, tag)
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
	colog.LDebug:   []byte("\u2699 "),
	colog.LInfo:    []byte("\u24d8 "),
	colog.LWarning: []byte("\u26a0 "),
	colog.LError:   []byte("\u2622 "),
	colog.LAlert:   []byte("\u2620 "),
}

var colorLabels = colog.LevelMap{
	colog.LTrace:   []byte("[ trace ] "),
	colog.LDebug:   []byte("\x1b[0;36m\u2699 \x1b[0m"),
	colog.LInfo:    []byte("\x1b[0;32m\u24d8 \x1b[0m"),
	colog.LWarning: []byte("\x1b[0;33m\u26a0 \x1b[0m"),
	colog.LError:   []byte("\x1b[0;31m\u2622 \x1b[0m"),
	colog.LAlert:   []byte("\x1b[0;37;41m\u2620 \x1b[0m"),
}

func initLog(ctx context.Context) error {
	lvl, err := colog.ParseLevel(ctx.LogLevel())
	if err != nil {
		return err
	}
	colog.SetMinLevel(lvl)
	colog.SetDefaultLevel(colog.LError)
	colog.SetFormatter(&colog.StdFormatter{
		Flag:        ctx.LogFlags(),
		HeaderPlain: plainLabels,
		HeaderColor: colorLabels,
	})
	colog.SetOutput(ctx.Stderr())
	colog.Register()
	return nil
}

func currentConfig(configFile string) (*context.Config, *context.Config, error) {
	var savedConfig *context.Config
	file, err := os.Open(configFile)
	switch {
	case err == nil:
		defer file.Close()
		savedConfig, err = context.LoadConfig(file)
		if err != nil {
			return nil, nil, err
		}
	case os.IsNotExist(err):
		savedConfig = &context.Config{}
	default:
		return nil, nil, err
	}

	savedConfig = context.MergeConfig(savedConfig, context.LoadKeyring())
	envarConfig, err := context.GetEnvarConfig()
	if err != nil {
		return nil, nil, err
	}
	cfg := context.MergeConfig(context.DefaultConfig(), savedConfig, envarConfig)
	if err := initLog(cfg); err != nil {
		return nil, nil, err
	}

	if err := context.ValidateContext(cfg); err != nil {
		log.Printf("warn: invalid config: %v", err)
	}
	return savedConfig, cfg, nil
}

func wrapCommand(cmd *kingpin.CmdClause, f func(context.Context) error) (string, func() error) {
	var configFile string
	setConfigFlag(cmd, &configFile)
	return cmd.FullCommand(), func() error {
		_, cfg, err := currentConfig(configFile)
		if err != nil {
			return err
		}

		return f(cfg)
	}
}

func wrapConfigurableCommand(cmd *kingpin.CmdClause, f func(*context.Config) error) (string, func() error) {
	var configFile string
	setConfigFlag(cmd, &configFile)
	return cmd.FullCommand(), func() error {
		savedConfig, _, err := currentConfig(configFile)
		if err != nil {
			return err
		}

		if err = f(savedConfig); err != nil {
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
		return context.SaveConfig(file, savedConfig)
	}
}
