package application

import (
	"context"
	"flag"
)

// type Application defines a common interface for command-line applications.
type Application interface {
	// DefaultFlagSet returns a `flag.FlagSet` instance with required flags and their default values assigned.
	DefaultFlagSet(context.Context) (*flag.FlagSet, error)
	// Run will run the code implementing the `Application` interface using a default flag set.
	Run(context.Context) error
	// Run will run the code implementing the `Application` interface using a custom flag set.
	RunWithFlagSet(context.Context, *flag.FlagSet) error
}
