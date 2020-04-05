package gordon

import (
	"strings"

	"github.com/blang/semver"
)

func ValidateTag(tag string) (semver.Version, error) {
	return semver.Parse(strings.TrimPrefix(tag, "v"))
}
