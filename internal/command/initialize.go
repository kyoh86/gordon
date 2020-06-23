package command

import (
	"context"
	"fmt"
)

func Initialize(_ context.Context, ev Env) error {
	fmt.Printf(`export PATH="%s:${PATH}"%s`, ev.Bin(), "\n")
	return nil
}
