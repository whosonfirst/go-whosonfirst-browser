package pruner

import (
)

type NullPruner struct {
	Pruner
}

func NewNullPruner() (Pruner, error) {

	pr := NullPruner{}
	return &pr, nil
}

func (pr *NullPruner) Prune(uri string) error {
     return nil
}
