package context

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/comail/colog"
	"github.com/pkg/errors"
)

var validOwnerRegexp = regexp.MustCompile(`^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$`)

func ValidateLogLevel(level string) error {
	_, err := colog.ParseLevel(level)
	return err
}

func ValidateOwner(owner string) error {
	if !validOwnerRegexp.MatchString(owner) {
		return errors.New("owner name may only contain alphanumeric characters or single hyphens, and cannot begin or end with a hyphen")
	}
	return nil
}

func ValidateFile(file string) error {
	path := filepath.Clean(file)
	info, err := os.Stat(path)
	switch {
	case err == nil:
		if info.IsDir() {
			return fmt.Errorf("is a directory: %s", file)
		}
		return nil
	case os.IsNotExist(err):
		return nil
	default:
		return err
	}
}

func ValidateRegexp(pattern string) error {
	_, err := regexp.Compile(pattern)
	return err
}

func ValidateRoot(root string) error {
	if root == "" {
		return errors.New("no root")
	}
	path := filepath.Clean(root)
	info, err := os.Stat(path)
	switch {
	case err == nil:
		if !info.IsDir() {
			return fmt.Errorf("is a file: %s", root)
		}
		return nil
	case os.IsNotExist(err):
		return nil
	default:
		return err
	}
}

func ValidateArchitecture(arch string) error {
	_, ok := goArchs[arch]
	if !ok {
		return errors.New("unsupported architecture")
	}

	return nil
}

func ValidateOS(os string) error {
	_, ok := goDists[os]
	if !ok {
		return errors.New("unsupported os")
	}

	return nil
}

func ValidateDist(os, arch string) error {
	goArchs, ok := goDists[os]
	if !ok {
		return errors.New("unsupported os")
	}

	_, ok = goArchs[arch]
	if !ok {
		return errors.New("unsupported architecture")
	}

	return nil
}

func ValidateContext(ctx Context) error {
	errs := map[string]error{}
	if err := ValidateLogLevel(ctx.LogLevel()); err != nil {
		errs["log.level"] = err
	}
	if err := ValidateRoot(ctx.Root()); err != nil {
		errs["root"] = err
	}
	if err := ValidateOwner(ctx.GitHubUser()); err != nil {
		errs["github.user"] = err
	}
	if err := ValidateDist(ctx.OS(), ctx.Architecture()); err != nil {
		errs["os|architecture"] = err
	}
	if len(errs) > 0 {
		sep := ""
		var builder strings.Builder
		builder.WriteString("invalid value {")
		for prop, err := range errs {
			builder.WriteString(prop)
			builder.WriteString(": ")
			builder.WriteString(err.Error())
			builder.WriteString(sep)
			sep = "; "
		}
		builder.WriteString("} ")
		return errors.New(builder.String())
	}
	return nil
}
