package filter

import (
	"net/url"
)

func NewSPRFilterFromQuery(query url.Values) (Filter, error) {

	inputs, err := NewSPRInputs()

	if err != nil {
		return nil, err
	}

	inputs.Placetypes = query["placetype"]
	inputs.IsCurrent = query["is_current"]
	inputs.IsDeprecated = query["is_deprecated"]
	inputs.IsCeased = query["is_ceased"]
	inputs.IsSuperseded = query["is_superseded"]
	inputs.IsSuperseding = query["is_superseding"]

	inputs.Geometries = query["geometries"]
	inputs.AlternateGeometries = query["alternate_geometry"]

	return NewSPRFilterFromInputs(inputs)
}
