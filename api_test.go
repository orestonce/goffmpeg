package goffmpeg

import (
	"fmt"
	"testing"
)

func TestSetupFfmpeg(t *testing.T) {
	p, err := SetupFfmpeg()
	fmt.Println(err, p)
	MustShowHelp(p)
}
