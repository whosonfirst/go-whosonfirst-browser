package update

// this might be moved in to its own go-whosonfirst-update package, but not today
// (20191223/thisisaaronland)

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/whosonfirst/go-whosonfirst-browser/schema"
	"log"
	"regexp"
)

type Update struct {
	Geometry   map[string]interface{}
	Properties map[string]interface{}
}

// this message signature will probably change (20191223/thisisaaronland)

func UpdateFeature(ctx context.Context, body []byte, update_req *Update, valid_paths *regexp.Regexp) ([]byte, int, error) {

	updates := 0

	var updated_body []byte
	var updated_err error

	for path, new_value := range update_req.Properties {

		path = fmt.Sprintf("properties.%s", path)

		log.Println("PATH", path)

		if !valid_paths.MatchString(path) {
			return nil, -1, errors.New("Invalid path")
		}

		// TO DO : SANITIZE value HERE... BUT WHAT ABOUT NOT-STRINGS...

		// TO DO: CHECK WHETHER PROPERTY/VALUE IS A SINGLETON - FOR EXAMPLE:
		// curl -s -X POST -F 'properties.wof:name=SPORK' -F 'properties.wof:name=BOB' http://localhost:8080/update/101736545 | \
		// python -mjson.tool | grep 'wof:name'
		// "wof:name": [

		// TO DO: HOW TO REMOVE THINGS...

		old_rsp := gjson.GetBytes(body, path)

		if old_rsp.Exists() {

			old_enc, err := json.Marshal(old_rsp.Value())

			if err != nil {
				return nil, -1, err
			}

			new_enc, err := json.Marshal(new_value)

			if err != nil {
				return nil, -1, err
			}

			if bytes.Compare(new_enc, old_enc) == 0 {
				log.Println("SAME SAME")
				continue
			}
		}

		def, err := schema.PropertiesDefinition(path)

		if err != nil {
			return nil, -1, err
		}

		log.Println("SET", path, new_value, def)

		updated_body, updated_err = sjson.SetBytes(body, path, new_value)

		if updated_err != nil {
			return nil, -1, updated_err
		}

		updates += 1
	}

	return updated_body, updates, nil
}
