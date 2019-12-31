package schema

import (
	"errors"
	"fmt"
	jsschema "github.com/lestrrat-go/jsschema"
	"github.com/lestrrat-go/jsschema/validator"
	"github.com/tidwall/gjson"
	"strings"
)

func GetPropertyDefinition(rel_path string) (string, error) {

	abs_path := fmt.Sprintf("definitions.properties.%s", rel_path)

	rsp := gjson.Get(PROPERTIES, abs_path)

	if !rsp.Exists() {
		return "", errors.New("Missing property definition.")
	}

	return rsp.String(), nil
}

func IsValidProperty(rel_path string, input interface{}) (bool, error) {

	def, err := GetPropertyDefinition(rel_path)

	if err != nil {
		return false, err
	}

	return Validate(def, input)
}

func HasValidGeometry(input interface{}) (bool, error) {

	rsp := gjson.Get(GEOMETRY, "definitions")

	if !rsp.Exists() {
		return false, errors.New("Missing geometry definition.")
	}

	return Validate(rsp.String(), input)
}

func HasValidProperties(input interface{}) (bool, error) {

	rsp := gjson.Get(PROPERTIES, "definitions.properties")

	if !rsp.Exists() {
		return false, errors.New("Missing properties definition.")
	}

	return Validate(rsp.String(), input)
}

func Validate(def string, input interface{}) (bool, error) {

	fh := strings.NewReader(def)
	s, err := jsschema.Read(fh)

	if err != nil {
		return false, err
	}

	v := validator.New(s)

	err = v.Validate(input)

	if err != nil {
		return false, err
	}

	return true, nil
}
