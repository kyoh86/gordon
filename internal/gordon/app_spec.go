package gordon

import (
	"errors"
	"flag"
	"fmt"
	"regexp"
	"strings"
)

// AppSpec specifies a release in the GitHub.
// AppSpec can accept "owner/name@tag" notation.
// If "@tag" is ommited, it means latest tag.
type AppSpec struct {
	raw   string
	owner string
	name  string
}

func (a AppSpec) Owner() string { return a.owner }
func (a AppSpec) Name() string  { return a.name }

// Set text as AppSpec
func (a *AppSpec) Set(rawApp string) error {
	terms := strings.Split(rawApp, "/")
	switch len(terms) {
	case 1:
		return fmt.Errorf("no slash in the app spec %q", rawApp)
	case 2:
		// noop
	default:
		return fmt.Errorf("plural slash in the app spec %q", rawApp)
	}

	if err := ValidateOwner(terms[0]); err != nil {
		return fmt.Errorf("invalid owner %w in the app spec %q", err, rawApp)
	}

	if err := ValidateName(terms[1]); err != nil {
		return fmt.Errorf("invalid name %w in the app spec %q", err, rawApp)
	}

	a.owner = terms[0]
	a.name = terms[1]
	a.raw = rawApp
	return nil
}

func (a AppSpec) String() string {
	return a.raw
}

var _ flag.Value = (*AppSpec)(nil)

func ParseAppSpec(rawApp string) (*AppSpec, error) {
	var app AppSpec
	if err := app.Set(rawApp); err != nil {
		return nil, err
	}
	return &app, nil
}

var invalidNameRegexp = regexp.MustCompile(`[^\w\-\.]`)

func ValidateName(name string) error {
	if name == "." {
		return errors.New("'.' is reserved name")
	}
	if name == ".." {
		return errors.New("'..' is reserved name")
	}
	if name == "" {
		return errors.New("project name is empty")
	}
	if invalidNameRegexp.MatchString(name) {
		return errors.New("invalid project name")
	}
	return nil
}

var validOwnerRegexp = regexp.MustCompile(`^[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*$`)

func ValidateOwner(owner string) error {
	if !validOwnerRegexp.MatchString(owner) {
		return errors.New("invalid owner name")
	}
	return nil
}
