package http

// maybe move this in to go-whosonfirst-uri proper?
// (20200116/thisisaaronland)

import (
	"errors"
	"github.com/whosonfirst/go-whosonfirst-uri"
	_ "log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var re_uri *regexp.Regexp

func init() {
	re_uri = regexp.MustCompile(`^(\d+)(?:\-alt(?:\-([a-zA-Z0-9_]+(?:\-[a-zA-Z0-9_]+(?:\-[a-zA-Z0-9_\-]+)?)?)))?(?:\.[^\.]+|\/)?$`)
}

/*

fmt.Printf("%q\n", re.FindStringSubmatch("12345.geojson"))
	fmt.Printf("%q\n", re.FindStringSubmatch("1234-alt-swim-path.geojson"))
	fmt.Printf("%q\n", re.FindStringSubmatch("1234-alt-swim-path-hello-world.geojson"))
	fmt.Printf("%q\n", re.FindStringSubmatch("1234-alt-swim.svg"))

*/

func IdFromURI(path string) (int64, *uri.URIArgs, error) {

	abs_path, err := filepath.Abs(path)

	if err != nil {
		return -1, nil, err
	}

	fname := filepath.Base(abs_path)

	match := re_uri.FindStringSubmatch(fname)

	// log.Println(fname, match)

	if len(match) == 0 {
		return -1, nil, errors.New("Unable to parse WOF ID")
	}

	if len(match) < 2 {
		return -1, nil, errors.New("Unable to parse WOF ID")
	}

	str_id := match[1]
	str_alt := match[2]

	wofid, err := strconv.ParseInt(str_id, 10, 64)

	if err != nil {
		return -1, nil, err
	}

	args := &uri.URIArgs{
		Alternate: false,
	}

	if str_alt != "" {

		args.Alternate = true

		alt := strings.Split(str_alt, "-")

		switch len(alt) {
		case 1:
			args.Source = alt[0]
		case 2:
			args.Source = alt[0]
			args.Function = alt[1]
		default:
			args.Source = alt[0]
			args.Function = alt[1]
			args.Extras = alt[2:]
		}
	}

	return wofid, args, nil
}
