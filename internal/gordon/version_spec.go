package gordon

import (
	"flag"
	"strings"
)

// VersionSpec specifies a release in the GitHub.
// VersionSpec can accept "owner/name@tag" notation.
// If "@tag" is ommited, it means latest tag.
type VersionSpec struct {
	raw string
	AppSpec
	tag string
}

func (v VersionSpec) Tag() string { return v.tag }

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

	v.tag = tag
	v.raw = rawRelease
	return nil
}

func (v VersionSpec) String() string {
	return v.raw
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
