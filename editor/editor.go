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
	_ "log"
	"regexp"
	"time"
)

type Editor struct {
	allowed_paths *regexp.Regexp
}

func NewEditor(allowed_paths *regexp.Regexp) (*Editor, error) {

	ed := &Editor{
		allowed_paths: allowed_paths,
	}

	return ed, nil
}

// this message signature may still change... (20191230/thisisaaronland)

func (ed *Editor) UpdateFeature(ctx context.Context, body []byte, update_req *UpdateRequest) ([]byte, *UpdateResponse, error) {

	updates := make([]*Update, 0)

	updated_body := make([]byte, len(body))
	copy(updated_body, body)

	var updated_err error

	if update_req.Geometry != nil {

		// SEE THIS? IT IS NOT POSSIBLE TO UPDATE JUST PART OF THE GEOMETRY
		// YOU HAVE TO UPDATE THE WHOLE THING (20191231/thisisaaronland)

		path := "geometry"

		if !ed.allowed_paths.MatchString(path) {
			return nil, nil, errors.New("Invalid path")
		}

		old_geom := gjson.GetBytes(body, "geometry")

		if !old_geom.Exists() {
			return nil, nil, errors.New("Missing geometry")
		}

		old_enc, err := json.Marshal(old_geom.Value())

		if err != nil {
			return nil, nil, err
		}

		new_enc, err := json.Marshal(update_req.Geometry)

		if err != nil {
			return nil, nil, err
		}

		if bytes.Compare(new_enc, old_enc) != 0 {

			// NOTE THIS IS _NOT_ TESTING WHETHER THE GEOMETRY/COORDINATES IS VALID

			_, err := schema.HasValidGeometry(update_req.Geometry)

			if err != nil {
				msg := fmt.Sprintf("geometry failed validation: %s", err.Error())
				return nil, nil, errors.New(msg)
			}

			// CHECK COORDINATES, WINDING, ETC. HERE...

			updated_body, updated_err = sjson.SetBytes(updated_body, path, update_req.Geometry)

			if updated_err != nil {
				return nil, nil, updated_err
			}

			u := &Update{
				Type: UPDATE_TYPE_CHANGE,
				Path: path,
			}

			updates = append(updates, u)
		}
	}

	for path, new_value := range update_req.Properties {

		path = fmt.Sprintf("properties.%s", path)

		if !ed.allowed_paths.MatchString(path) {
			return nil, nil, errors.New("Invalid path")
		}

		// TBD - is this the best way to signal things to delete? as
		// { "properties": { "foo.bar.baz": null } }
		// (20191230/thisisaaronland)

		if new_value == nil {

			updated_body, updated_err = sjson.DeleteBytes(updated_body, path)

			if updated_err != nil {
				return nil, nil, updated_err
			}

			u := &Update{
				Type: UPDATE_TYPE_REMOVE,
				Path: path,
			}

			updates = append(updates, u)
			continue
		}

		// TO DO : SANITIZE value HERE... BUT WHAT ABOUT NOT-STRINGS...

		old_rsp := gjson.GetBytes(body, path)

		if old_rsp.Exists() {

			old_enc, err := json.Marshal(old_rsp.Value())

			if err != nil {
				return nil, nil, err
			}

			new_enc, err := json.Marshal(new_value)

			if err != nil {
				return nil, nil, err
			}

			if bytes.Compare(new_enc, old_enc) == 0 {
				continue
			}
		}

		_, err := schema.IsValidProperty(path, new_value)

		if err != nil {
			msg := fmt.Sprintf("'%s' property failed validation: %s", path, err.Error())
			return nil, nil, errors.New(msg)
		}

		// log.Println("SET", path, new_value)

		updated_body, updated_err = sjson.SetBytes(updated_body, path, new_value)

		if updated_err != nil {
			return nil, nil, updated_err
		}

		u := &Update{
			Type: UPDATE_TYPE_CHANGE,
			Path: path,
		}

		updates = append(updates, u)
	}

	update_rsp := &UpdateResponse{
		Updates: updates,
	}

	return updated_body, update_rsp, nil
}

// something something something EDTF dates... (20191229/straup)

func (ed *Editor) DeprecateFeature(ctx context.Context, body []byte, t time.Time) ([]byte, *UpdateResponse, error) {

	deprecated_rsp := gjson.GetBytes(body, "properties.edtf:deprecated")

	if deprecated_rsp.Exists() {
		return nil, nil, errors.New("Feature is already deprecated")
	}

	props := map[string]interface{}{
		"edtf:deprecated": t.Format("2006-01-02"),
		"mz:is_current":   0,
	}

	update_req := &UpdateRequest{
		Properties: props,
	}

	return ed.UpdateFeature(ctx, body, update_req)
}

// something something something EDTF dates... (20191229/straup)

func (ed *Editor) CessateFeature(ctx context.Context, body []byte, t time.Time) ([]byte, *UpdateResponse, error) {

	cessated_rsp := gjson.GetBytes(body, "properties.edtf:cessated")

	if cessated_rsp.Exists() {
		return nil, nil, errors.New("Feature is already cessated")
	}

	props := map[string]interface{}{
		"edtf:cessation": t.Format("2006-01-02"),
		"mz:is_current":  0,
	}

	update_req := &UpdateRequest{
		Properties: props,
	}

	return ed.UpdateFeature(ctx, body, update_req)
}
