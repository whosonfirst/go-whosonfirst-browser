package application

import (
	"context"
	"flag"
)

type Application interface {
	DefaultFlagSet(context.Context) (*flag.FlagSet, error)
	Run(context.Context) error
	RunWithFlagSet(context.Context, *flag.FlagSet) error
}
