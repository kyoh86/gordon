// Code generated by main.go DO NOT EDIT.

package env

import (
	"fmt"
	types "github.com/kyoh86/appenv/types"
	"io"
)

type Config struct {
	yml YAML
}

func GetConfig(yamlReader io.Reader) (config Config, err error) {
	yml, err := loadYAML(yamlReader)
	if err != nil {
		return config, err
	}
	return buildConfig(yml)
}

func buildConfig(yml YAML) (config Config, err error) {
	config.yml = yml
	return
}

func (c *Config) Save(yamlWriter io.Writer) error {
	if err := saveYAML(yamlWriter, &c.yml); err != nil {
		return err
	}
	return nil
}

func PropertyNames() []string {
	return []string{"github.host", "github.user", "architecture", "os", "cache", "bin", "man", "root", "hooks"}
}

func (a *Config) Property(name string) (types.Config, error) {
	switch name {
	case "github.host":
		return &githubHostConfig{parent: a}, nil
	case "github.user":
		return &githubUserConfig{parent: a}, nil
	case "architecture":
		return &architectureConfig{parent: a}, nil
	case "os":
		return &osConfig{parent: a}, nil
	case "cache":
		return &cacheConfig{parent: a}, nil
	case "bin":
		return &binConfig{parent: a}, nil
	case "man":
		return &manConfig{parent: a}, nil
	case "root":
		return &rootConfig{parent: a}, nil
	case "hooks":
		return &hooksConfig{parent: a}, nil
	}
	return nil, fmt.Errorf("invalid property name %q", name)
}

type githubHostConfig struct {
	parent *Config
}

func (a *githubHostConfig) Get() (string, error) {
	{
		p := a.parent.yml.GithubHost
		if p != nil {
			text, err := p.MarshalText()
			return string(text), err
		}
	}
	return "", nil
}

func (a *githubHostConfig) Set(value string) error {
	{
		p := a.parent.yml.GithubHost
		if p == nil {
			p = new(GithubHost)
		}
		if err := p.UnmarshalText([]byte(value)); err != nil {
			return err
		}
		a.parent.yml.GithubHost = p
	}
	return nil
}

func (a *githubHostConfig) Unset() {
	a.parent.yml.GithubHost = nil
}

type githubUserConfig struct {
	parent *Config
}

func (a *githubUserConfig) Get() (string, error) {
	{
		p := a.parent.yml.GithubUser
		if p != nil {
			text, err := p.MarshalText()
			return string(text), err
		}
	}
	return "", nil
}

func (a *githubUserConfig) Set(value string) error {
	{
		p := a.parent.yml.GithubUser
		if p == nil {
			p = new(GithubUser)
		}
		if err := p.UnmarshalText([]byte(value)); err != nil {
			return err
		}
		a.parent.yml.GithubUser = p
	}
	return nil
}

func (a *githubUserConfig) Unset() {
	a.parent.yml.GithubUser = nil
}

type architectureConfig struct {
	parent *Config
}

func (a *architectureConfig) Get() (string, error) {
	{
		p := a.parent.yml.Architecture
		if p != nil {
			text, err := p.MarshalText()
			return string(text), err
		}
	}
	return "", nil
}

func (a *architectureConfig) Set(value string) error {
	{
		p := a.parent.yml.Architecture
		if p == nil {
			p = new(Architecture)
		}
		if err := p.UnmarshalText([]byte(value)); err != nil {
			return err
		}
		a.parent.yml.Architecture = p
	}
	return nil
}

func (a *architectureConfig) Unset() {
	a.parent.yml.Architecture = nil
}

type osConfig struct {
	parent *Config
}

func (a *osConfig) Get() (string, error) {
	{
		p := a.parent.yml.OS
		if p != nil {
			text, err := p.MarshalText()
			return string(text), err
		}
	}
	return "", nil
}

func (a *osConfig) Set(value string) error {
	{
		p := a.parent.yml.OS
		if p == nil {
			p = new(OS)
		}
		if err := p.UnmarshalText([]byte(value)); err != nil {
			return err
		}
		a.parent.yml.OS = p
	}
	return nil
}

func (a *osConfig) Unset() {
	a.parent.yml.OS = nil
}

type cacheConfig struct {
	parent *Config
}

func (a *cacheConfig) Get() (string, error) {
	{
		p := a.parent.yml.Cache
		if p != nil {
			text, err := p.MarshalText()
			return string(text), err
		}
	}
	return "", nil
}

func (a *cacheConfig) Set(value string) error {
	{
		p := a.parent.yml.Cache
		if p == nil {
			p = new(Cache)
		}
		if err := p.UnmarshalText([]byte(value)); err != nil {
			return err
		}
		a.parent.yml.Cache = p
	}
	return nil
}

func (a *cacheConfig) Unset() {
	a.parent.yml.Cache = nil
}

type binConfig struct {
	parent *Config
}

func (a *binConfig) Get() (string, error) {
	{
		p := a.parent.yml.Bin
		if p != nil {
			text, err := p.MarshalText()
			return string(text), err
		}
	}
	return "", nil
}

func (a *binConfig) Set(value string) error {
	{
		p := a.parent.yml.Bin
		if p == nil {
			p = new(Bin)
		}
		if err := p.UnmarshalText([]byte(value)); err != nil {
			return err
		}
		a.parent.yml.Bin = p
	}
	return nil
}

func (a *binConfig) Unset() {
	a.parent.yml.Bin = nil
}

type manConfig struct {
	parent *Config
}

func (a *manConfig) Get() (string, error) {
	{
		p := a.parent.yml.Man
		if p != nil {
			text, err := p.MarshalText()
			return string(text), err
		}
	}
	return "", nil
}

func (a *manConfig) Set(value string) error {
	{
		p := a.parent.yml.Man
		if p == nil {
			p = new(Man)
		}
		if err := p.UnmarshalText([]byte(value)); err != nil {
			return err
		}
		a.parent.yml.Man = p
	}
	return nil
}

func (a *manConfig) Unset() {
	a.parent.yml.Man = nil
}

type rootConfig struct {
	parent *Config
}

func (a *rootConfig) Get() (string, error) {
	{
		p := a.parent.yml.Root
		if p != nil {
			text, err := p.MarshalText()
			return string(text), err
		}
	}
	return "", nil
}

func (a *rootConfig) Set(value string) error {
	{
		p := a.parent.yml.Root
		if p == nil {
			p = new(Root)
		}
		if err := p.UnmarshalText([]byte(value)); err != nil {
			return err
		}
		a.parent.yml.Root = p
	}
	return nil
}

func (a *rootConfig) Unset() {
	a.parent.yml.Root = nil
}

type hooksConfig struct {
	parent *Config
}

func (a *hooksConfig) Get() (string, error) {
	{
		p := a.parent.yml.Hooks
		if p != nil {
			text, err := p.MarshalText()
			return string(text), err
		}
	}
	return "", nil
}

func (a *hooksConfig) Set(value string) error {
	{
		p := a.parent.yml.Hooks
		if p == nil {
			p = new(Hooks)
		}
		if err := p.UnmarshalText([]byte(value)); err != nil {
			return err
		}
		a.parent.yml.Hooks = p
	}
	return nil
}

func (a *hooksConfig) Unset() {
	a.parent.yml.Hooks = nil
}