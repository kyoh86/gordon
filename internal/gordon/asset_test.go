package gordon

import "testing"

func TestIsAlnum(t *testing.T) {
	if !isAlnum('a') {
		t.Error("'a' is alnum, but...")
	}
	if !isAlnum('z') {
		t.Error("'z' is alnum, but...")
	}
	if !isAlnum('A') {
		t.Error("'A' is alnum, but...")
	}
	if !isAlnum('Z') {
		t.Error("'Z' is alnum, but...")
	}
	if !isAlnum('0') {
		t.Error("'0' is alnum, but...")
	}
	if !isAlnum('9') {
		t.Error("'9' is alnum, but...")
	}
	if isAlnum('-') {
		t.Error("'-' is not alnum, but...")
	}
	if isAlnum('_') {
		t.Error("'_' is not alnum, but...")
	}
}

func TestContainsWord(t *testing.T) {
	containsSet := []struct {
		title  string
		input  string
		substr string
	}{{
		title:  "whole string",
		input:  "whole",
		substr: "whole",
	}, {
		title:  "underscored",
		input:  "_sub_str_",
		substr: "sub_str",
	}, {
		title:  "wrapped other words",
		input:  "whole_sub_str_remain",
		substr: "sub_str",
	}, {
		title:  "starts with",
		input:  "substr_foo",
		substr: "substr",
	}, {
		title:  "ends with",
		input:  "foo_substr",
		substr: "substr",
	}}
	for _, contains := range containsSet {
		contains := contains
		t.Run(contains.title, func(t *testing.T) {
			if !containsWord(contains.input, contains.substr) {
				t.Errorf("substr %q is not found in %q", contains.substr, contains.input)
			}
		})
	}

	notContainsSet := []struct {
		title  string
		input  string
		substr string
	}{{
		title:  "empty string",
		input:  "",
		substr: "substr",
	}, {
		title:  "overwrapped string",
		input:  "sub",
		substr: "substr",
	}, {
		title:  "following alnum rune",
		input:  "0substr",
		substr: "substr",
	}, {
		title:  "followed by alnum rune",
		input:  "substrZ",
		substr: "substr",
	}, {
		title:  "following alnum string",
		input:  "01substr",
		substr: "substr",
	}, {
		title:  "followed by alnum string",
		input:  "substrYZ",
		substr: "substr",
	}}
	for _, notContains := range notContainsSet {
		notContains := notContains
		t.Run(notContains.title, func(t *testing.T) {
			if containsWord(notContains.input, notContains.substr) {
				t.Errorf("substr %q is found in %q", notContains.substr, notContains.input)
			}
		})
	}
}

func TestMatchArchitecture(t *testing.T) {
	if !MatchArchitecture("foo_0.3.4_mac_64bit", "amd64") {
		t.Errorf(`expect "amd64" matches "foo_0.3.4_mac_64bit", but not`)
	}
}
