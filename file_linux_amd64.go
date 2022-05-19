package goffmpeg

import (
	_ "embed"
)

//go:embed w_linux-ia32.gz
var gEmbedFfmpegGz []byte
