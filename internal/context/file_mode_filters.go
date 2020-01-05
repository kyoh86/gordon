package context

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type FileMode os.FileMode

func (f FileMode) Match(m os.FileMode) bool {
	return (os.FileMode(f) & m) == os.FileMode(f)
}

func (f FileMode) MarshalText() ([]byte, error) {
	buf := bytes.Buffer{}
	if _, err := fmt.Fprintf(&buf, "%04o", f); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (f *FileMode) UnmarshalText(raw []byte) error {
	if len(raw) != 4 {
		return errors.New("invalid format")
	}
	u, err := strconv.ParseUint(string(raw), 8, 32)
	if err != nil {
		return err
	}
	*f = FileMode(os.FileMode(u))
	return nil
}

type FileModes []FileMode

const FileModeSeparator = '|'

func (filters FileModes) MarshalText() ([]byte, error) {
	if len(filters) == 0 {
		return nil, nil
	}
	result := make([]byte, 0, 5*len(filters)-1)
	first, err := filters[0].MarshalText()
	if err != nil {
		return nil, err
	}
	result = append(result, first...)
	for _, m := range filters[1:] {
		result = append(result, byte(FileModeSeparator))
		concat, err := m.MarshalText()
		if err != nil {
			return nil, err
		}
		result = append(result, concat...)
	}
	return result, nil
}

func (filters *FileModes) UnmarshalText(raw []byte) error {
	if len(raw) == 0 {
		*filters = nil
		return nil
	}
	words := bytes.Split(raw, []byte{byte(FileModeSeparator)})
	result := make(FileModes, 0, len(words))
	for _, word := range words {
		var m FileMode
		if err := m.UnmarshalText(word); err != nil {
			return err
		}
		result = append(result, m)
	}
	*filters = result
	return nil
}

func (filters FileModes) Match(m os.FileMode) bool {
	for _, filter := range filters {
		if filter.Match(m) {
			return true
		}
	}
	return false
}
