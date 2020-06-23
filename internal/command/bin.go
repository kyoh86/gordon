package command

import (
	"context"
	"fmt"
)

func Bin(_ context.Context, ev Env) error {
	fmt.Println(ev.Bin())
	return nil
}
