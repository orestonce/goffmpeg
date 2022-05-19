package goffmpeg

import (
	_ "embed"
)

//go:embed w_win32-ia32.gz
var gEmbedFfmpegGz []byte
