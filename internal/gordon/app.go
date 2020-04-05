package gordon

import (
	"errors"
	"os"
)

type App struct {
	owner string
	name  string
}

var (
	ErrAppNotFound = errors.New("app not found")
)

func FindApp(ev Env, spec AppSpec) (*App, error) {
	fi, err := os.Stat(assetSubPath(ev, spec.owner, spec.name))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrAppNotFound
		}
		return nil, err
	}
	if !fi.IsDir() {
		return nil, ErrAppNotFound
	}
	return &App{
		owner: spec.owner,
		name:  spec.name,
	}, nil
}
