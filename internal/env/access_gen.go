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

	access.historyFile = new(HistoryFile).Default().(string)
	if yml.HistoryFile != nil {
		access.historyFile = yml.HistoryFile.Value().(string)
	}
	if envar.HistoryFile != nil {
		access.historyFile = envar.HistoryFile.Value().(string)
	}

	access.historySave = new(HistorySave).Default().(bool)
	if yml.HistorySave != nil {
		access.historySave = yml.HistorySave.Value().(bool)
	}
	if envar.HistorySave != nil {
		access.historySave = envar.HistorySave.Value().(bool)
	}

	access.extractModes = new(ExtractModes).Default().([]os.FileMode)
	if yml.ExtractModes != nil {
		access.extractModes = yml.ExtractModes.Value().([]os.FileMode)
	}
	if envar.ExtractModes != nil {
		access.extractModes = envar.ExtractModes.Value().([]os.FileMode)
	}

	access.extractExclude = new(ExtractExclude).Default().(string)
	if yml.ExtractExclude != nil {
		access.extractExclude = yml.ExtractExclude.Value().(string)
	}
	if envar.ExtractExclude != nil {
		access.extractExclude = envar.ExtractExclude.Value().(string)
	}

	access.extractInclude = new(ExtractInclude).Default().(string)
	if yml.ExtractInclude != nil {
		access.extractInclude = yml.ExtractInclude.Value().(string)
	}
	if envar.ExtractInclude != nil {
		access.extractInclude = envar.ExtractInclude.Value().(string)
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
	githubHost     string
	githubUser     string
	historyFile    string
	historySave    bool
	extractModes   []os.FileMode
	extractExclude string
	extractInclude string
	architecture   string
	os             string
	root           string
	hooks          []string
}

func (a *Access) GithubHost() string {
	return a.githubHost
}

func (a *Access) GithubUser() string {
	return a.githubUser
}

func (a *Access) HistoryFile() string {
	return a.historyFile
}

func (a *Access) HistorySave() bool {
	return a.historySave
}

func (a *Access) ExtractModes() []os.FileMode {
	return a.extractModes
}

func (a *Access) ExtractExclude() string {
	return a.extractExclude
}

func (a *Access) ExtractInclude() string {
	return a.extractInclude
}

func (a *Access) Architecture() string {
	return a.architecture
}

func (a *Access) OS() string {
	return a.os
}

func (a *Access) Root() string {
	return a.root
}

func (a *Access) Hooks() []string {
	return a.hooks
}
