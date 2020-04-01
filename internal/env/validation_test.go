package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateOwner(t *testing.T) {
	expect := "owner name may only contain alphanumeric characters or single hyphens, and cannot begin or end with a hyphen"
	assert.EqualError(t, ValidateOwner(""), expect, "fail when empty owner is given")
	assert.EqualError(t, ValidateOwner("kyoh_86"), expect, "fail when owner name contains invalid character")
	assert.EqualError(t, ValidateOwner("-kyoh86"), expect, "fail when owner name starts with hyphen")
	assert.EqualError(t, ValidateOwner("kyoh86-"), expect, "fail when owner name ends with hyphen")
	assert.NoError(t, ValidateOwner("kyoh86"), "success")
}

func TestValidateRoot(t *testing.T) {
	assert.EqualError(t, ValidateRoot(""), "no root", "fail when no path in root")
	assert.NoError(t, ValidateRoot("/path/to/not/existing"))
	assert.Error(t, ValidateRoot("\x00"))
}

func TestValidateContext(t *testing.T) {
	t.Run("invalid root", func(t *testing.T) {
		ctx := &MockContext{
			MRoot:       "/\x00",
			MLogLevel:   "warn",
			MGitHubUser: "kyoh86",
		}
		assert.Error(t, ValidateContext(ctx))
	})
	t.Run("invalid owner", func(t *testing.T) {
		ctx := &MockContext{
			MRoot:       "/path/to/not/existing",
			MLogLevel:   "warn",
			MGitHubUser: "",
		}
		assert.Error(t, ValidateContext(ctx))
	})
	t.Run("invalid loglevel", func(t *testing.T) {
		ctx := &MockContext{
			MRoot:       "/path/to/not/existing",
			MLogLevel:   "invalid",
			MGitHubUser: "kyoh86",
		}
		assert.Error(t, ValidateContext(ctx))
	})
	t.Run("valid env", func(t *testing.T) {
		ctx := &MockContext{
			MRoot:         "/path/to/not/existing",
			MLogLevel:     "warn",
			MGitHubUser:   "kyoh86",
			MOS:           "linux",
			MArchitecture: "386",
		}
		assert.NoError(t, ValidateContext(ctx))
	})
}
