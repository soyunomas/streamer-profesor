package web

import "embed"

//go:embed static/index.html
var Assets embed.FS
