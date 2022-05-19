package goffmpeg

import (
	_ "embed"
)

//go:embed w_darwin-arm64.gz
var gEmbedFfmpegGz []byte

