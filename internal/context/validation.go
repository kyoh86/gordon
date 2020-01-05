package context

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/comail/colog"
	"github.com/pkg/errors"
)

var invalidNameRegexp = regexp.MustCompile(`[^\w\-\.]`)

func ValidateName(name string) error {
	if name == "." || name == ".." {
		return errors.New("'.' or '..' is reserved name")
	}
	if name == "" {
		return errors.New("empty project name")
	}
	if invalidNameRegexp.MatchString(name) {
		return errors.New("project name may only contain alphanumeric characters, dots or hyphens")
	}
	return nil
}

var validOwnerRegexp = regexp.MustCompile(`^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$`)

func ValidateOwner(owner string) error {
	if !validOwnerRegexp.MatchString(owner) {
		return errors.New("owner name may only contain alphanumeric characters or single hyphens, and cannot begin or end with a hyphen")
	}
	return nil
}

func ValidateRoot(root string) error {
	if root == "" {
		return errors.New("no root")
	}
	path := filepath.Clean(root)
	_, err := os.Stat(path)
	switch {
	case err == nil:
		return nil
	case os.IsNotExist(err):
		return nil
	default:
		return err
	}
}

func ValidateLogLevel(level string) error {
	_, err := colog.ParseLevel(level)
	return err
}

func ValidateContext(ctx Context) error {
	if err := ValidateRoot(ctx.Root()); err != nil {
		return errors.Wrap(err, "invalid root in the context; set a valid path by 'gordon config put root <Project root path>'")
	}
	if err := ValidateOwner(ctx.GitHubUser()); err != nil {
		return errors.Wrap(err, "invalid GitHub user in the context; set a valid name by 'gordon config put github.user <GitHub user name>'")
	}
	if err := ValidateLogLevel(ctx.LogLevel()); err != nil {
		return errors.Wrap(err, "invalid log level in the context; set a valid log-level by 'gordon config put log.level <debug|info|warn|error>' or unset by 'gordon config unset log.level'")
	}
	return nil
}
