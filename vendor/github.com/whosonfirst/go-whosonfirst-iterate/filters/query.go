package filters

import (
	"context"
	"github.com/aaronland/go-json-query"
	"io"
	"net/url"
)

type QueryFilters struct {
	Filters
	Include *query.QuerySet
	Exclude *query.QuerySet
}

func NewQueryFiltersFromURI(ctx context.Context, uri string) (Filters, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()
	return NewQueryFiltersFromQuery(ctx, q)
}

func NewQueryFiltersFromQuery(ctx context.Context, q url.Values) (Filters, error) {

	f := &QueryFilters{}

	includes := q["include"]
	excludes := q["exclude"]

	if len(includes) > 0 {

		includes_mode := q.Get("include_mode")
		includes_qs, err := querySetFromStrings(ctx, includes_mode, includes...)

		if err != nil {
			return nil, err
		}

		f.Include = includes_qs
	}

	if len(excludes) > 0 {

		excludes_mode := q.Get("exclude_mode")
		excludes_qs, err := querySetFromStrings(ctx, excludes_mode, excludes...)

		if err != nil {
			return nil, err
		}

		f.Exclude = excludes_qs
	}

	return f, nil
}

func querySetFromStrings(ctx context.Context, mode string, flags ...string) (*query.QuerySet, error) {

	var query_flags query.QueryFlags
	query_mode := query.QUERYSET_MODE_ALL

	for _, fl := range flags {

		err := query_flags.Set(fl)

		if err != nil {
			return nil, err
		}
	}

	switch mode {
	case query.QUERYSET_MODE_ALL, query.QUERYSET_MODE_ANY:
		query_mode = mode
	default:
		// pass
	}

	qs := &query.QuerySet{
		Queries: query_flags,
		Mode:    query_mode,
	}

	return qs, nil
}

func (f *QueryFilters) Apply(ctx context.Context, fh io.ReadSeekCloser) (bool, error) {

	body, err := io.ReadAll(fh)

	if err != nil {
		return false, err
	}

	includes_qs := f.Include
	excludes_qs := f.Exclude

	if includes_qs != nil {

		matches, err := query.Matches(ctx, includes_qs, body)

		if err != nil {
			return false, err
		}

		if !matches {
			return false, nil
		}

	}

	if excludes_qs != nil {

		matches, err := query.Matches(ctx, excludes_qs, body)

		if err != nil {
			return false, err
		}

		if matches {
			return false, nil
		}
	}

	return true, nil
}
