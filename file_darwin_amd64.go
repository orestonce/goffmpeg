package goffmpeg

import (
	_ "embed"
)

//go:embed w_darwin-x64.gz
var gEmbedFfmpegGz []byte
