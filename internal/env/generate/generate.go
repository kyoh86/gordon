// +build generate

package main

import (
	"log"

	"github.com/kyoh86/appenv"
	"github.com/kyoh86/gordon/internal/env"
)

//go:generate go run -tags generate .

func main() {
	g := &appenv.Generator{}

	if err := g.Do(
		"github.com/kyoh86/gordon/internal/env",
		"../",
		appenv.Opt(new(env.GithubHost), appenv.StoreYAML(), appenv.StoreEnvar()),
		appenv.Opt(new(env.GithubUser), appenv.StoreYAML(), appenv.StoreEnvar()),
		appenv.Opt(new(env.Architecture), appenv.StoreYAML(), appenv.StoreEnvar()),
		appenv.Opt(new(env.OS), appenv.StoreYAML(), appenv.StoreEnvar()),
		appenv.Opt(new(env.Cache), appenv.StoreYAML(), appenv.StoreEnvar()),
		appenv.Opt(new(env.Bin), appenv.StoreYAML(), appenv.StoreEnvar()),
		appenv.Opt(new(env.Man), appenv.StoreYAML(), appenv.StoreEnvar()),
		appenv.Opt(new(env.Hooks), appenv.StoreYAML(), appenv.StoreEnvar()),
	); err != nil {
		log.Fatalln(err)
	}
}
