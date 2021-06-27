package static

import (
	"embed"
)

//go:embed css/* javascript/* fonts/*
var FS embed.FS
