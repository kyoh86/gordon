package gordon

import (
	"errors"
	"flag"
	"fmt"
	"regexp"
	"strings"
)

// AppSpec specifies a release in the GitHub.
// AppSpec can accept "owner/name@version" notation.
// If "@version" is ommited, it means latest version.
type AppSpec struct {
	raw     string
	owner   string
	name    string
	version string
}

func (a AppSpec) Owner() string   { return a.owner }
func (a AppSpec) Name() string    { return a.name }
func (a AppSpec) Version() string { return a.version }

// Set text as AppSpec
func (a *AppSpec) Set(rawApp string) error {
	var version string

	terms := strings.SplitN(rawApp, "@", 2)
	if len(terms) == 2 {
		version = terms[1]
	}

	terms = strings.Split(terms[0], "/")
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
	a.version = version
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

// AppSpecs is array of AppSpec
type AppSpecs []AppSpec

// Set will add a text to AppSpecs as a AppSpec
func (a *AppSpecs) Set(value string) error {
	app := new(AppSpec)
	if err := app.Set(value); err != nil {
		return err
	}
	*a = append(*a, *app)
	return nil
}

func (a AppSpecs) String() string {
	if len(a) == 0 {
		return ""
	}
	strs := make([]string, 0, len(a))
	for _, app := range a {
		strs = append(strs, app.String())
	}
	return strings.Join(strs, ",")
}

func (AppSpecs) IsCumulative() bool { return true }

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
