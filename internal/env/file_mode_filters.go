package env

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"

	types "github.com/kyoh86/appenv/types"
)

type FileMode struct {
	value os.FileMode
}

func (f FileMode) Match(m os.FileMode) bool {
	return matchFileMode(f.value, m)
}

func matchFileMode(filter, mode os.FileMode) bool {
	return (filter & mode) == filter
}

func (f *FileMode) Value() interface{} {
	return f.value
}

func (*FileMode) Default() interface{} {
	return os.FileMode(0)
}

func (f FileMode) MarshalText() ([]byte, error) {
	return marshalFileModeText(f.value)
}

func marshalFileModeText(f os.FileMode) ([]byte, error) {
	buf := bytes.Buffer{}
	if _, err := fmt.Fprintf(&buf, "%04o", f); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (f *FileMode) UnmarshalText(raw []byte) error {
	v, err := unmarshalFileModeText(raw)
	if err != nil {
		return err
	}

	f.value = v
	return nil
}

func unmarshalFileModeText(raw []byte) (os.FileMode, error) {
	if len(raw) != 4 {
		return 0, errors.New("invalid format")
	}
	u, err := strconv.ParseUint(string(raw), 8, 32)
	if err != nil {
		return 0, err
	}
	return os.FileMode(u), nil
}

var _ types.Value = (*FileMode)(nil)

type FileModes struct {
	values []os.FileMode
}

func (f *FileModes) Value() interface{} {
	return f.values
}

func (*FileModes) Default() interface{} {
	return nil
}

const FileModeSeparator = '|'

func (f FileModes) MarshalText() ([]byte, error) {
	if len(f.values) == 0 {
		return nil, nil
	}
	result := make([]byte, 0, 5*len(f.values)-1)
	first, err := marshalFileModeText(f.values[0])
	if err != nil {
		return nil, err
	}
	result = append(result, first...)
	for _, m := range f.values[1:] {
		result = append(result, byte(FileModeSeparator))
		concat, err := marshalFileModeText(m)
		if err != nil {
			return nil, err
		}
		result = append(result, concat...)
	}
	return result, nil
}

func (f *FileModes) UnmarshalText(raw []byte) error {
	if len(raw) == 0 {
		f.values = nil
		return nil
	}
	words := bytes.Split(raw, []byte{byte(FileModeSeparator)})
	result := make([]os.FileMode, 0, len(words))
	for _, word := range words {
		m, err := unmarshalFileModeText(word)
		if err != nil {
			return err
		}
		result = append(result, m)
	}
	f.values = result
	return nil
}

func (f FileModes) Match(m os.FileMode) bool {
	for _, filter := range f.values {
		if matchFileMode(filter, m) {
			return true
		}
	}
	return false
}

var _ types.Value = (*FileModes)(nil)
