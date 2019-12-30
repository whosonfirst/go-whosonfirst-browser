package editor

import (
	"fmt"
	"strings"
)

const UPDATE_TYPE_CHANGE string = "change"
const UPDATE_TYPE_REMOVE string = "remove"

type UpdateRequest struct {
	Geometry   map[string]interface{}
	Properties map[string]interface{}
}

// TBD: return Feature/GeoJSON body as part of UpdateResponse
// (20191203/thisisaaronland)

type UpdateResponse struct {
	Updates []*Update `json:"updates"`
}

func (u *UpdateResponse) String() string {

	if u.Count() == 0 {
		return "No updates."
	}

	str_updates := make([]string, len(u.Updates))

	for i, u := range u.Updates {
		str_updates[i] = fmt.Sprintf("* %s", u.String())
	}

	return strings.Join(str_updates, "\n")
}

func (u *UpdateResponse) Count() int {
	return len(u.Updates)
}

type Update struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

func (u *Update) String() string {
	return fmt.Sprintf("[%s] %s", u.Type, u.Path)
}
