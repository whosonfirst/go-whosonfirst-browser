package pruner

import (

)

type Pruner interface {
	Prune(string) error
}
