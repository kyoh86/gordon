package hub

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type TokenManager interface {
	SetGithubToken(user, token string) error
	GetGithubToken(user string) (string, error)
	DeleteGithubToken(user string) error
}

type MemoryTokenManager struct {
	m map[string]string
}

func NewMemory(s string) (TokenManager, error) {
	return &MemoryTokenManager{m: map[string]string{}}, nil
}

func (m *MemoryTokenManager) SetGithubToken(user, token string) error {
	m.m[user] = token
	return nil
}

func (m *MemoryTokenManager) GetGithubToken(user string) (string, error) {
	return m.m[user], nil
}

func (m *MemoryTokenManager) DeleteGithubToken(user string) error {
	delete(m.m, user)
	return nil
}

type FileTokenManager struct {
	file string
}

func NewFileForHost(host string) (TokenManager, error) {
	cache, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(cache, "gordon")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	file := filepath.Join(dir, host+"_tokens.json")
	return NewFile(file)
}

func NewFile(file string) (TokenManager, error) {
	return &FileTokenManager{file: file}, nil
}

func readFile(file string) (map[string]string, error) {
	fp, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	defer fp.Close()
	var m map[string]string
	if err := json.NewDecoder(fp).Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}

func writeFile(file string, m map[string]string) error {
	fp, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()
	return json.NewEncoder(fp).Encode(m)
}

func (m *FileTokenManager) SetGithubToken(user, token string) error {
	obj, err := readFile(m.file)
	if err != nil {
		return err
	}
	obj[user] = token
	return writeFile(m.file, obj)
}

func (m *FileTokenManager) GetGithubToken(user string) (string, error) {
	obj, err := readFile(m.file)
	if err != nil {
		return "", err
	}
	return obj[user], nil
}

func (m *FileTokenManager) DeleteGithubToken(user string) error {
	obj, err := readFile(m.file)
	if err != nil {
		return err
	}
	delete(obj, user)
	return writeFile(m.file, obj)
}
