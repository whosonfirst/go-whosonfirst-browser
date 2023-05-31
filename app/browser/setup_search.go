package browser

import (
	"context"
	"fmt"
	"sync"

	"github.com/whosonfirst/go-whosonfirst-search/fulltext"
)

var setupSearchOnce sync.Once
var setupSearchError error

func setupSearch() {

	ctx := context.Background()
	var err error

	search_database, err = fulltext.NewFullTextDatabase(ctx, cfg.SearchDatabaseURI)

	if err != nil {
		setupSearchError = fmt.Errorf("Failed to create fulltext database for '%s', %w", cfg.SearchDatabaseURI, err)
		return
	}

}
