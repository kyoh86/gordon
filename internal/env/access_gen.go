// Code generated by main.go DO NOT EDIT.

package env

import "io"

func GetAccess(yamlReader io.Reader, envarPrefix string) (access Access, err error) {
	yml, err := loadYAML(yamlReader)
	if err != nil {
		return access, err
	}
	return buildAccess(yml, envarPrefix)
}

func buildAccess(yml YAML, envarPrefix string) (access Access, err error) {
	envar, err := getEnvar(envarPrefix)
	if err != nil {
		return access, err
	}
	access.githubHost = new(GithubHost).Default().(string)
	if yml.GithubHost != nil {
		access.githubHost = yml.GithubHost.Value().(string)
	}
	if envar.GithubHost != nil {
		access.githubHost = envar.GithubHost.Value().(string)
	}

	access.githubUser = new(GithubUser).Default().(string)
	if yml.GithubUser != nil {
		access.githubUser = yml.GithubUser.Value().(string)
	}
	if envar.GithubUser != nil {
		access.githubUser = envar.GithubUser.Value().(string)
	}

	access.architecture = new(Architecture).Default().(string)
	if yml.Architecture != nil {
		access.architecture = yml.Architecture.Value().(string)
	}
	if envar.Architecture != nil {
		access.architecture = envar.Architecture.Value().(string)
	}

	access.os = new(OS).Default().(string)
	if yml.OS != nil {
		access.os = yml.OS.Value().(string)
	}
	if envar.OS != nil {
		access.os = envar.OS.Value().(string)
	}

	access.cache = new(Cache).Default().(string)
	if yml.Cache != nil {
		access.cache = yml.Cache.Value().(string)
	}
	if envar.Cache != nil {
		access.cache = envar.Cache.Value().(string)
	}

	access.bin = new(Bin).Default().(string)
	if yml.Bin != nil {
		access.bin = yml.Bin.Value().(string)
	}
	if envar.Bin != nil {
		access.bin = envar.Bin.Value().(string)
	}

	access.man = new(Man).Default().(string)
	if yml.Man != nil {
		access.man = yml.Man.Value().(string)
	}
	if envar.Man != nil {
		access.man = envar.Man.Value().(string)
	}

	access.root = new(Root).Default().(string)
	if yml.Root != nil {
		access.root = yml.Root.Value().(string)
	}
	if envar.Root != nil {
		access.root = envar.Root.Value().(string)
	}

	access.hooks = new(Hooks).Default().([]string)
	if yml.Hooks != nil {
		access.hooks = yml.Hooks.Value().([]string)
	}
	if envar.Hooks != nil {
		access.hooks = envar.Hooks.Value().([]string)
	}

	return
}

type Access struct {
	githubHost   string
	githubUser   string
	architecture string
	os           string
	cache        string
	bin          string
	man          string
	root         string
	hooks        []string
}

func (a *Access) GithubHost() string {
	return a.githubHost
}

func (a *Access) GithubUser() string {
	return a.githubUser
}

func (a *Access) Architecture() string {
	return a.architecture
}

func (a *Access) OS() string {
	return a.os
}

func (a *Access) Cache() string {
	return a.cache
}

func (a *Access) Bin() string {
	return a.bin
}

func (a *Access) Man() string {
	return a.man
}

func (a *Access) Root() string {
	return a.root
}

func (a *Access) Hooks() []string {
	return a.hooks
}