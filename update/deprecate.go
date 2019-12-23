package update

import (
	"context"
	"errors"
	"github.com/tidwall/gjson"
	"regexp"
	"time"
)

func DeprecateFeature(ctx context.Context, body []byte, valid_paths *regexp.Regexp) ([]byte, error) {

	deprecated_rsp := gjson.GetBytes(body, "properties.edtf:deprecated")

	if deprecated_rsp.Exists() {
		return nil, errors.New("Feature is already deprecated")
	}

	now := time.Now()

	props := map[string]interface{}{
		"edtf:deprecated": now.Format("2006-01-02"),
		"mz:is_current":   0,
	}

	update_req := &Update{
		Properties: props,
	}

	updated_body, _, err := UpdateFeature(ctx, body, update_req, valid_paths)
	return updated_body, err
}
