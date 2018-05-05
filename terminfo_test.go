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

	if len(ti.Names) != 2 {
		t.Errorf("Expected 2 names, got %d\n", len(ti.Names))
	}

	if ti.Names[0] != "xterm-256color" || ti.Names[1] != "xterm with 256 colors" {
		t.Errorf("Name mismatch\n")
	}
}

func TestConstants(t *testing.T) {
	if BoolLPIChangesRes != 36 {
		t.Error("Boolean constant mismatch")
	}
	if NumBitImageType != 32 {
		t.Error("Number constant mismatch")
	}
	if StrSetPglenInch != 393 {
		t.Error("String constant mismatch")
	}
}
