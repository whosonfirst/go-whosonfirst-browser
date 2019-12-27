package schema

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
)

func PropertiesDefinition(rel_path string) (*gjson.Result, error) {

	abs_path := fmt.Sprintf("definitions.properties.%s", rel_path)

	rsp := gjson.Get(PROPERTIES, abs_path)

	if !rsp.Exists() {
		return nil, errors.New("Undefined property")
	}

	return &rsp, nil
}
