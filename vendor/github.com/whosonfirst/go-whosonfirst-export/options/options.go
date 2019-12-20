package options

import (
	"github.com/whosonfirst/go-whosonfirst-export/uid"
)

type Options interface {
	UIDProvider() uid.Provider

	// mmmmmmmmaybe? (20190110/thisisaaronland)
	// Get(string) (interface{}, bool)
	// Set(string, interface{}) error
}
