package application

import (
	"context"
	"flag"
)

// type Application defines a common interface for command-line applications.
type Application interface {
	DefaultFlagSet(context.Context) (*flag.FlagSet, error)
	Run(context.Context) error
	RunWithFlagSet(context.Context, *flag.FlagSet) error
}
