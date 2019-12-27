package update

import (
	"context"
	"errors"
	"github.com/tidwall/gjson"
	"regexp"
	"time"
)

func CessateFeature(ctx context.Context, body []byte, valid_paths *regexp.Regexp) ([]byte, error) {

	cessated_rsp := gjson.GetBytes(body, "properties.edtf:cessation")

	if cessated_rsp.Exists() && cessated_rsp.String() != "uuuu" {
		return nil, errors.New("Feature is already cessated")
	}

	now := time.Now()

	props := map[string]interface{}{
		"edtf:cessation": now.Format("2006-01-02"),
		"mz:is_current":  0,
	}

	update_req := &Update{
		Properties: props,
	}

	updated_body, _, err := UpdateFeature(ctx, body, update_req, valid_paths)
	return updated_body, err
}
