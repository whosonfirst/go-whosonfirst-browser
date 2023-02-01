# go-log

Opinionated Go package for doing minimal structured logging and prefixing of log messages with Emoji for easier filtering.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/aaronland/go-log.svg)](https://pkg.go.dev/github.com/aaronland/go-log)

## Example

```
import (
	"log"

	aa_log "github.com/aaronland/go-log"
)	

func main(){

	logger := log.Default()

	aa_log.SetMinLevelWithPrefix(aa_log.WARNING_PREFIX)

	// No output
	aa_log.Debug(logger, "This is a test")

	aa_log.UnsetMinLevel()

	// prints "🪵 This is a second test"
	aa_log.Debug(logger, "This is a second test")
}
```

## Prefixes

| Log level | Prefix |
| --- | --- |
| debug | 🪵 |
| info | 💬 |
| warning | 🧯 |
| error | 🔥 |
| fatal | 💥 |
