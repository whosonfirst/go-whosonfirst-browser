package properties

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"sync"
)

func EnsureBelongsTo(feature []byte) ([]byte, error) {

	belongsto := make([]int64, 0)

	rsp := gjson.GetBytes(feature, "properties.wof:hierarchy")

	if rsp.Exists() {

		s := new(sync.Map)

		for _, h := range rsp.Array() {

			for _, i := range h.Map() {
				id := i.Int()
				s.Store(id, true)
			}
		}

		s.Range(func(k interface{}, v interface{}) bool {
			id := k.(int64)

			if id > 0 {
				belongsto = append(belongsto, id)
			}

			return true
		})
	}

	return sjson.SetBytes(feature, "properties.wof:belongsto", belongsto)
}
