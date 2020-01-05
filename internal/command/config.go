package command

import (
	"fmt"

	"github.com/kyoh86/gordon/internal/context"
)

func ConfigGetAll(cfg *context.Config) error {
	for _, name := range context.OptionNames() {
		opt, _ := context.Option(name) // ignore error: context.OptionNames covers all accessor
		value := opt.Get(cfg)
		fmt.Printf("%s = %s\n", name, value)
	}
	return nil
}

func ConfigGet(cfg *context.Config, optionName string) error {
	opt, err := context.Option(optionName)
	if err != nil {
		return err
	}
	value := opt.Get(cfg)
	fmt.Println(value)
	return nil
}

func ConfigPut(cfg *context.Config, optionName, optionValue string) error {
	opt, err := context.Option(optionName)
	if err != nil {
		return err
	}
	return opt.Put(cfg, optionValue)
}

func ConfigUnset(cfg *context.Config, optionName string) error {
	opt, err := context.Option(optionName)
	if err != nil {
		return err
	}
	return opt.Unset(cfg)
}
