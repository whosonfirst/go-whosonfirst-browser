package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-index"
	"github.com/whosonfirst/go-whosonfirst-index/utils"
	"github.com/whosonfirst/go-whosonfirst-log"
	"github.com/whosonfirst/go-whosonfirst-mysql"
	"github.com/whosonfirst/go-whosonfirst-mysql/database"
	"github.com/whosonfirst/go-whosonfirst-mysql/tables"
	"github.com/whosonfirst/warning"
	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func main() {

	valid_modes := strings.Join(index.Modes(), ",")
	desc_modes := fmt.Sprintf("The mode to use importing data. Valid modes are: %s.", valid_modes)

	config := flag.String("config", "", "Read some or all flags from an ini-style config file. Values in the config file take precedence over command line flags.")
	section := flag.String("section", "wof-mysql", "A valid ini-style config file section.")

	dsn := flag.String("dsn", "", "A valid go-sql-driver DSN string, for example '{USER}:{PASSWORD}@/{DATABASE}'")
	mode := flag.String("mode", "repo", desc_modes)

	index_geojson := flag.Bool("geojson", false, "Index the 'geojson' tables")
	index_whosonfirst := flag.Bool("whosonfirst", false, "Index the 'whosonfirst' tables")
	index_all := flag.Bool("all", false, "Index all the tables")

	timings := flag.Bool("timings", false, "Display timings during and after indexing")

	flag.Parse()

	logger := log.SimpleWOFLogger()

	stdout := io.Writer(os.Stdout)
	logger.AddLogger(stdout, "status")

	if *config != "" {

		cfg, err := ini.LoadSources(ini.LoadOptions{
			AllowBooleanKeys: true,
		}, *config)

		if err != nil {
			logger.Fatal("Unable to load config file because %s", err)
		}

		sect, err := cfg.GetSection(*section)

		if err != nil {
			logger.Fatal("Config file is missing 'mysql' section, %s", err)
		}

		flag.VisitAll(func(fl *flag.Flag) {

			name := fl.Name

			if name == "section" {
				return
			}

			if sect.HasKey(name) {

				k := sect.Key(name)
				v := k.Value()

				flag.Set(name, v)

				logger.Status("Reset %s flag from config file", name)
			}
		})

	} else {

		flag.VisitAll(func(fl *flag.Flag) {

			name := fl.Name
			env_name := name

			env_name = strings.Replace(name, "-", "_", -1)
			env_name = strings.ToUpper(name)

			env_name = fmt.Sprintf("WOF_MYSQL_%s", env_name)

			v, ok := os.LookupEnv(env_name)

			if ok {

				flag.Set(name, v)
				logger.Status("Reset %s flag from %s environment variable", name, env_name)
			}
		})
	}

	db, err := database.NewDB(*dsn)

	if err != nil {
		logger.Fatal("unable to create database (%s) because %s", *dsn, err)
	}

	defer db.Close()

	to_index := make([]mysql.Table, 0)

	if *index_geojson || *index_all {

		tbl, err := tables.NewGeoJSONTableWithDatabase(db)

		if err != nil {
			logger.Fatal("failed to create 'geojson' table because %s", err)
		}

		to_index = append(to_index, tbl)
	}

	if *index_whosonfirst || *index_all {

		tbl, err := tables.NewWhosonfirstTableWithDatabase(db)

		if err != nil {
			logger.Fatal("failed to create 'whosonfirst' table because %s", err)
		}

		to_index = append(to_index, tbl)
	}

	if len(to_index) == 0 {
		logger.Fatal("You forgot to specify which (any) tables to index")
	}

	table_timings := make(map[string]time.Duration)
	mu := new(sync.RWMutex)

	cb := func(fh io.Reader, ctx context.Context, args ...interface{}) error {

		path, err := index.PathForContext(ctx)

		if err != nil {
			return err
		}

		ok, err := utils.IsPrincipalWOFRecord(fh, ctx)

		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		f, err := feature.LoadWOFFeatureFromReader(fh)

		if err != nil {

			if err != nil && !warning.IsWarning(err) {
				msg := fmt.Sprintf("Unable to load %s, because %s", path, err)
				return errors.New(msg)
			}

		}

		db.Lock()

		defer db.Unlock()

		for _, t := range to_index {

			t1 := time.Now()

			err = t.IndexFeature(db, f)

			if err != nil {
				logger.Warning("failed to index feature (%s) in '%s' table because %s", path, t.Name(), err)
				return nil
			}

			t2 := time.Since(t1)

			n := t.Name()

			mu.Lock()

			_, ok := table_timings[n]

			if ok {
				table_timings[n] += t2
			} else {
				table_timings[n] = t2
			}

			mu.Unlock()
		}

		return nil
	}

	indexer, err := index.NewIndexer(*mode, cb)

	if err != nil {
		logger.Fatal("Failed to create new indexer because: %s", err)
	}

	done_ch := make(chan bool)
	t1 := time.Now()

	show_timings := func() {

		t2 := time.Since(t1)

		i := atomic.LoadInt64(&indexer.Indexed) // please just make this part of go-whosonfirst-index

		mu.RLock()
		defer mu.RUnlock()

		for t, d := range table_timings {
			logger.Status("time to index %s (%d) : %v", t, i, d)
		}

		logger.Status("time to index all (%d) : %v", i, t2)
	}

	if *timings {

		go func() {

			for {

				select {
				case <-done_ch:
					return
				case <-time.After(1 * time.Minute):
					show_timings()
				}
			}
		}()

		defer func() {
			done_ch <- true
		}()
	}

	err = indexer.IndexPaths(flag.Args())

	if err != nil {
		logger.Fatal("Failed to index paths in %s mode because: %s", *mode, err)
	}

	os.Exit(0)
}
