package gordon

import (
	"flag"
	"strings"

	"github.com/blang/semver"
)

// VersionSpec specifies a release in the GitHub.
// VersionSpec can accept "owner/name@tag" notation.
// If "@tag" is ommited, it means latest tag.
type VersionSpec struct {
	AppSpec
	tag    string
	semver semver.Version
}

func (v VersionSpec) Tag() string { return v.tag }

func (v VersionSpec) WithoutTag() VersionSpec {
	return VersionSpec{
		AppSpec: v.AppSpec,
	}
}

// Set text as VersionSpec
func (v *VersionSpec) Set(rawRelease string) error {
	var tag string

	terms := strings.SplitN(rawRelease, "@", 2)
	if len(terms) == 2 {
		tag = terms[1]
	}

	if err := v.AppSpec.Set(terms[0]); err != nil {
		return err
	}

	if tag != "" {
		sv, err := ValidateTag(tag)
		if err != nil {
			return err
		}
		v.tag = tag
		v.semver = sv
	}
	return nil
}

func (v VersionSpec) String() string {
	if v.tag == "" {
		return v.AppSpec.String()
	}
	return v.AppSpec.String() + "@" + v.tag
}

var _ flag.Value = (*VersionSpec)(nil)

func ParseVersionSpec(rawRelease string) (*VersionSpec, error) {
	var ver VersionSpec
	if err := ver.Set(rawRelease); err != nil {
		return nil, err
	}
	return &ver, nil
}

// VersionSpecs is array of VersionSpec
type VersionSpecs []VersionSpec

// Set will add a text to VersionSpecs as a VersionSpec
func (v *VersionSpecs) Set(value string) error {
	ver := new(VersionSpec)
	if err := ver.Set(value); err != nil {
		return err
	}
	*v = append(*v, *ver)
	return nil
}

func (v VersionSpecs) String() string {
	if len(v) == 0 {
		return ""
	}
	strs := make([]string, 0, len(v))
	for _, ver := range v {
		strs = append(strs, ver.String())
	}
	return strings.Join(strs, ",")
}

func (VersionSpecs) IsCumulative() bool { return true }
