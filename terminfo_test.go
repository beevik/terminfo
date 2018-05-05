package terminfo

import (
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	file, err := os.Open("/usr/share/terminfo/x/xterm1")
	if err != nil {
		t.Errorf("%v\n", err)
	}
	defer file.Close()

	ti, err := Read(file)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	_ = ti
}
