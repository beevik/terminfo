package terminfo

import (
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	file, err := os.Open("./_data/xterm-256color")
	if err != nil {
		t.Errorf("%v\n", err)
	}
	defer file.Close()

	ti, err := Read(file)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	names := ti.Names()
	if len(names) != 2 {
		t.Errorf("Expected 2 names, got %d\n", len(names))
	}

	if names[0] != "xterm-256color" || names[1] != "xterm with 256 colors" {
		t.Errorf("Name mismatch\n")
	}
}
