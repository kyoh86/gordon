package env

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestFileMode(t *testing.T) {
	t.Run("Marshal/UnmarshalText", func(t *testing.T) {
		var m FileMode
		assert.Error(t, m.UnmarshalText([]byte("1q84")))
		assert.Equal(t, FileMode(0), m)
		assert.NoError(t, m.UnmarshalText([]byte("0111")))
		assert.Equal(t, FileMode(0111), m)
		str, err := m.MarshalText()
		assert.NoError(t, err)
		assert.Equal(t, "0111", string(str))
	})
	t.Run("Marshal/UnmarshalJSON", func(t *testing.T) {
		var obj struct {
			Mode FileMode `json:"mode"`
		}
		assert.Error(t, json.Unmarshal([]byte(`{"mode":"1q84"}`), &obj))
		assert.Equal(t, FileMode(0), obj.Mode)
		assert.NoError(t, json.Unmarshal([]byte(`{"mode":"0111"}`), &obj))
		assert.Equal(t, FileMode(0111), obj.Mode)
		str, err := json.Marshal(obj)
		assert.NoError(t, err)
		assert.Equal(t, `{"mode":"0111"}`, string(str))
	})
	t.Run("Marshal/UnmarshalYAML", func(t *testing.T) {
		var obj struct {
			Mode FileMode `yaml:"mode"`
		}
		assert.Error(t, yaml.Unmarshal([]byte(`mode: "1q84"`+"\n"), &obj))
		assert.Equal(t, FileMode(0), obj.Mode)
		assert.NoError(t, yaml.Unmarshal([]byte(`mode: "0111"`+"\n"), &obj))
		assert.Equal(t, FileMode(0111), obj.Mode)
		str, err := yaml.Marshal(obj)
		assert.NoError(t, err)
		assert.Equal(t, `mode: "0111"`+"\n", string(str))
	})
}

func TestFileModes(t *testing.T) {
	t.Run("Marshal/UnmarshalText", func(t *testing.T) {
		var m FileModes
		assert.Error(t, m.UnmarshalText([]byte("0111|1q84")))
		assert.Empty(t, m)
		assert.NoError(t, m.UnmarshalText([]byte("0111|0222")))
		if assert.Len(t, m, 2) {
			assert.Equal(t, FileMode(0111), m[0])
			assert.Equal(t, FileMode(0222), m[1])
		}
		str, err := m.MarshalText()
		assert.NoError(t, err)
		assert.Equal(t, "0111|0222", string(str))
	})
	t.Run("Marshal/UnmarshalJSON", func(t *testing.T) {
		var obj struct {
			Mode FileModes `json:"mode"`
		}
		assert.Error(t, json.Unmarshal([]byte(`{"mode":"0111|1q84"}`), &obj))
		assert.Empty(t, obj.Mode)
		assert.NoError(t, json.Unmarshal([]byte(`{"mode":"0111|0222"}`), &obj))
		if assert.Len(t, obj.Mode, 2) {
			assert.Equal(t, FileMode(0111), obj.Mode[0])
			assert.Equal(t, FileMode(0222), obj.Mode[1])
		}
		str, err := json.Marshal(obj)
		assert.NoError(t, err)
		assert.Equal(t, `{"mode":"0111|0222"}`, string(str))
	})
	t.Run("Marshal/UnmarshalYAML", func(t *testing.T) {
		var obj struct {
			Mode FileModes `yaml:"mode"`
		}
		assert.Error(t, yaml.Unmarshal([]byte(`mode: 0111|1q84`), &obj))
		assert.Empty(t, obj.Mode)
		assert.NoError(t, yaml.Unmarshal([]byte(`mode: 0111`), &obj))
		if assert.Len(t, obj.Mode, 1) {
			assert.Equal(t, FileMode(0111), obj.Mode[0])
		}
		assert.NoError(t, yaml.Unmarshal([]byte(`mode: 0111|0222`), &obj))
		if assert.Len(t, obj.Mode, 2) {
			assert.Equal(t, FileMode(0111), obj.Mode[0])
			assert.Equal(t, FileMode(0222), obj.Mode[1])
		}
		str, err := yaml.Marshal(obj)
		assert.NoError(t, err)
		assert.Equal(t, `mode: 0111|0222`+"\n", string(str))
	})
}
