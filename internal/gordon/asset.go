package gordon

import "strings"

var archAliases = map[string][]string{
	"386":   {"i386"},
	"amd64": {"x86_64", "86_64"},
	"arm":   {"arm32"},
}

func isAlnum(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || '0' <= r && r <= '9'
}

func containsWord(s, substr string) bool {
	s = strings.ToLower(s)
	runes := []rune(s)
	for i := 0; i < len(s); {
		index := strings.Index(s[i:], substr)
		if index < 0 {
			return false
		}
		if (index <= 0 || !isAlnum(runes[index-1])) &&
			(len(s)-1 < index+len(substr) || !isAlnum(runes[index+len(substr)])) {
			return true
		}
		i = index + 1
	}
	return false
}

func MatchArchitecture(s, architecture string) bool {
	if containsWord(s, architecture) {
		return true
	}
	for _, alias := range archAliases[architecture] {
		if containsWord(s, alias) {
			return true
		}
	}
	return false
}

var osAliases = map[string][]string{
	"darwin":  {"osx", "mac", "macintosh"},
	"windows": {"win"},
}

func MatchOS(s, os string) bool {
	if containsWord(s, os) {
		return true
	}
	for _, alias := range osAliases[os] {
		if containsWord(s, alias) {
			return true
		}
	}
	return false
}
