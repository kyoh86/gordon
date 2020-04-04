package gordon_test

import (
	"strings"
	"testing"

	"github.com/kyoh86/gordon/internal/gordon"
)

func TestAppSpec(t *testing.T) {
	var specs gordon.AppSpecs
	if !specs.IsCumulative() {
		t.Fatalf("kingpin.Flag accepts multi values if the `IsCumulative` returns true")
	}
	if specs.String() != "" {
		t.Errorf("empty list should be empty string, but %q", specs.String())
	}
	for _, success := range []struct {
		title         string
		input         string
		expectOwner   string
		expectName    string
		expectVersion string
	}{{
		title:         "full-spec",
		input:         "kyoh86/gordon@v0.0.1",
		expectOwner:   "kyoh86",
		expectName:    "gordon",
		expectVersion: "v0.0.1",
	}, {
		title:         "no-version",
		input:         "kyoh86/gordon",
		expectOwner:   "kyoh86",
		expectName:    "gordon",
		expectVersion: "",
	}, {
		title:         "complex-version",
		input:         "kyoh86/gordon@v0.0.1-alpha@jp/us",
		expectOwner:   "kyoh86",
		expectName:    "gordon",
		expectVersion: "v0.0.1-alpha@jp/us",
	}} {
		t.Run(success.title, func(t *testing.T) {
			app, err := gordon.ParseAppSpec(success.input)
			if err != nil {
				t.Fatal(err)
			}
			if success.expectOwner != app.Owner() {
				t.Errorf("expect Owner  %q but %q", success.expectOwner, app.Owner())
			}
			if success.expectName != app.Name() {
				t.Errorf("expect Name   %q but %q", success.expectName, app.Name())
			}
			if success.expectVersion != app.Version() {
				t.Errorf("expect Version %q but %q", success.expectVersion, app.Version())
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
			app, err := gordon.ParseAppSpec(failure.input)
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
		"kyoh86/gordon@v0.0.1-alpha@jp/us",
	}, ",")
	result := specs.String()
	if joined != result {
		t.Errorf("specs.String() should be raws separated with comma, but %q", result)
	}
}
