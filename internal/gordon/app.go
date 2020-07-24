package gordon

import (
	"errors"
	"fmt"
)

type App struct {
	owner string
	name  string
}

func (a App) Spec() AppSpec {
	return AppSpec{
		owner: a.owner,
		name:  a.name,
	}
}

func (a App) String() string {
	return fmt.Sprintf("%s/%s", a.owner, a.name)
}

var (
	ErrAppNotFound = errors.New("app not found")
)

func FindApp(ev Env, spec AppSpec) (*App, error) {
	ok, err := isDir(assetSubPath(ev, spec.owner, spec.name))
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrAppNotFound
	}
	return &App{
		owner: spec.owner,
		name:  spec.name,
	}, nil
}
