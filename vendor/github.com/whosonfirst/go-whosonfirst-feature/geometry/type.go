package geometry

import (
	"fmt"
	"github.com/tidwall/gjson"
)

func Type(body []byte) (string, error) {

	rsp := gjson.GetBytes(body, "geometry.type")

	if !rsp.Exists() {
		return "", fmt.Errorf("Missing geometry.type property")
	}

	return rsp.String(), nil
}
