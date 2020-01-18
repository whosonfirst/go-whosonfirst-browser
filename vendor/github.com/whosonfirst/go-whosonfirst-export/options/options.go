package options

import (
	"github.com/whosonfirst/go-whosonfirst-id"
)

type Options interface {
	IDProvider() id.Provider
}
