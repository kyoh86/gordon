package context

import (
	"errors"
	"fmt"
	"strings"
)

type BoolOption string

var (
	TrueOption      = BoolOption("yes")
	FalseOption     = BoolOption("no")
	EmptyBoolOption = BoolOption("")
)

func (c BoolOption) String() string {
	return string(c)
}

func (c BoolOption) Bool() bool {
	return c == TrueOption
}

func (c *BoolOption) SetBool(b bool) {
	if b {
		*c = TrueOption
	} else {
		*c = FalseOption
	}
}

// Decode implements the interface `envdecode.Decoder`
func (c *BoolOption) Decode(repl string) error {
	switch strings.ToLower(repl) {
	case "yes", "no", "":
		*c = BoolOption(repl)
		return nil
	}
	return errors.New("invalid type")
}

// MarshalText implements the interface `encoding.TextMarshaler`
func (c BoolOption) MarshalText() ([]byte, error) {
	return []byte(c), nil
}

// UnmarshalText implements the interface `encoding.TextUnmarshaler`
func (c *BoolOption) UnmarshalText(raw []byte) error {
	switch string(raw) {
	case "yes":
		*c = TrueOption
	case "no":
		*c = FalseOption
	case "":
		*c = EmptyBoolOption
	default:
		return fmt.Errorf("invalid bool option %#v", string(raw))
	}
	return nil
}
