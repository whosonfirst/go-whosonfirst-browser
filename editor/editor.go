package editor

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
	"time"
)

type Editor struct {
	allowed_paths *regexp.Regexp
}

type Update struct {
	Geometry   map[string]interface{}
	Properties map[string]interface{}
}

func NewEditor(allowed_paths *regexp.Regexp) (*Editor, error) {

	ed := &Editor{
		allowed_paths: allowed_paths,
	}

	return ed, nil
}

// this message signature will probably change (20191223/thisisaaronland)

func (ed *Editor) UpdateFeature(ctx context.Context, body []byte, update_req *Update) ([]byte, int, error) {

	updates := 0

	updated_body := make([]byte, len(body))
	copy(updated_body, body)

	var updated_err error

	for path, new_value := range update_req.Properties {

		path = fmt.Sprintf("properties.%s", path)

		if !ed.allowed_paths.MatchString(path) {
			return nil, -1, errors.New("Invalid path")
		}

		// TO DO : SANITIZE value HERE... BUT WHAT ABOUT NOT-STRINGS...

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
				continue
			}
		}

		_, err := schema.IsValidProperty(path, new_value)

		if err != nil {
			msg := fmt.Sprintf("'%s' property failed validation: %s", path, err.Error())
			return nil, -1, errors.New(msg)
		}

		log.Println("SET", path, new_value)

		updated_body, updated_err = sjson.SetBytes(updated_body, path, new_value)

		if updated_err != nil {
			return nil, -1, updated_err
		}

		updates += 1
	}

	return updated_body, updates, nil
}

// something something something EDTF dates... (20191229/straup)

func (ed *Editor) DeprecateFeature(ctx context.Context, body []byte, t time.Time) ([]byte, error) {

	deprecated_rsp := gjson.GetBytes(body, "properties.edtf:deprecated")

	if deprecated_rsp.Exists() {
		return nil, errors.New("Feature is already deprecated")
	}

	props := map[string]interface{}{
		"edtf:deprecated": t.Format("2006-01-02"),
		"mz:is_current":   0,
	}

	update_req := &Update{
		Properties: props,
	}

	updated_body, _, err := ed.UpdateFeature(ctx, body, update_req)
	return updated_body, err
}

// something something something EDTF dates... (20191229/straup)

func (ed *Editor) CessateFeature(ctx context.Context, body []byte, t time.Time) ([]byte, error) {

	cessated_rsp := gjson.GetBytes(body, "properties.edtf:cessated")

	if cessated_rsp.Exists() {
		return nil, errors.New("Feature is already cessated")
	}

	props := map[string]interface{}{
		"edtf:cessation": t.Format("2006-01-02"),
		"mz:is_current":  0,
	}

	update_req := &Update{
		Properties: props,
	}

	updated_body, _, err := ed.UpdateFeature(ctx, body, update_req)
	return updated_body, err
}
