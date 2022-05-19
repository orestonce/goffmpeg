package goffmpeg

import (
	_ "embed"
)

//go:embed w_linux-arm.gz
var gEmbedFfmpegGz []byte
