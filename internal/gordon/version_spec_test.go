package gordon_test

import (
	"strings"
	"testing"

	"github.com/kyoh86/gordon/internal/gordon"
)

func TestVersionSpec(t *testing.T) {
	var specs gordon.VersionSpecs
	if !specs.IsCumulative() {
		t.Fatalf("kingpin.Flag accepts multi values if the `IsCumulative` returns true")
	}
	if specs.String() != "" {
		t.Errorf("empty list should be empty string, but %q", specs.String())
	}
	for _, success := range []struct {
		title       string
		input       string
		expectOwner string
		expectName  string
		expectTag   string
	}{{
		title:       "full-spec",
		input:       "kyoh86/gordon@v0.0.1",
		expectOwner: "kyoh86",
		expectName:  "gordon",
		expectTag:   "v0.0.1",
	}, {
		title:       "no-tag",
		input:       "kyoh86/gordon",
		expectOwner: "kyoh86",
		expectName:  "gordon",
		expectTag:   "",
	}, {
		title:       "complex-tag",
		input:       "kyoh86/gordon@v0.0.1+112.us-alpha.jp.us",
		expectOwner: "kyoh86",
		expectName:  "gordon",
		expectTag:   "v0.0.1+112.us-alpha.jp.us",
	}} {
		t.Run(success.title, func(t *testing.T) {
			app, err := gordon.ParseVersionSpec(success.input)
			if err != nil {
				t.Fatal(err)
			}
			if success.expectOwner != app.Owner() {
				t.Errorf("expect Owner  %q but %q", success.expectOwner, app.Owner())
			}
			if success.expectName != app.Name() {
				t.Errorf("expect Name   %q but %q", success.expectName, app.Name())
			}
			if success.expectTag != app.Tag() {
				t.Errorf("expect Tag %q but %q", success.expectTag, app.Tag())
			}
			if success.input != app.String() {
				t.Errorf("expect String %q but %q", success.input, app.String())
			}
			if err := specs.Set(success.input); err != nil {
				t.Errorf("failed to put app spec into the list %q", err)
			}
		})
	}

	for _, failure := range []struct {
		title string
		input string
	}{{
		title: "empty",
		input: "",
	}, {
		title: "no slash",
		input: "gordon",
	}, {
		title: "empty owner",
		input: "/gordon",
	}, {
		title: "invalid char in owner",
		input: "kyoh_86/gordon",
	}, {
		title: "invalid start in owner",
		input: "-kyoh86/gordon",
	}, {
		title: "invalid end in owner",
		input: "kyoh86-/gordon",
	}, {
		title: "empty name",
		input: "kyoh86/",
	}, {
		title: "single dot name",
		input: "kyoh86/.",
	}, {
		title: "invalid char in name",
		input: "kyoh86/gordon:thomas",
	}, {
		title: "double dot name",
		input: "kyoh86/..",
	}, {
		title: "too many slashes",
		input: "kyoh86/gordon/thomas",
	}} {
		t.Run(failure.title, func(t *testing.T) {
			app, err := gordon.ParseVersionSpec(failure.input)
			if err == nil {
				t.Errorf("expect error but nil and got app %q", app)
			}
			if err := specs.Set(failure.input); err == nil {
				t.Errorf("expect to fail to put app spec into the list but not")
			}
		})
	}
	joined := strings.Join([]string{
		"kyoh86/gordon@v0.0.1",
		"kyoh86/gordon",
		"kyoh86/gordon@v0.0.1+112.us-alpha.jp.us",
	}, ",")
	result := specs.String()
	if joined != result {
		t.Errorf("specs.String() should be raws separated with comma, but %q", result)
	}
}
