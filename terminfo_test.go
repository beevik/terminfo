package terminfo

import (
	"os"
	"testing"
)

type testCase struct {
	filename string
	names    []string
	boolCaps map[string]bool
	numCaps  map[string]int
	strCaps  map[string]string
}

var testCases = []testCase{
	{
		filename: "./_data/xterm-256color",
		names:    []string{"xterm-256color", "xterm with 256 colors"},
		boolCaps: map[string]bool{"xhp": false, "hc": false, "npc": true},
		numCaps:  map[string]int{"cols": 80, "it": 8, "lines": 24, "colors": 256, "pairs": 32767},
		strCaps:  map[string]string{"kf5": "\x1b[15~", "kRIT5": "\x1b[1;5C", "cub": "\x1b[%p1%dD"},
	},
	{
		filename: "./_data/rxvt",
		names:    []string{"rxvt", "rxvt terminal emulator (X Window System)"},
		boolCaps: map[string]bool{"xhp": false, "hc": false, "npc": false},
		numCaps:  map[string]int{"cols": 80, "it": 8, "lines": 24, "colors": 8, "pairs": 64},
		strCaps:  map[string]string{"kf5": "\x1b[15~", "kRIT5": "\x1bOc", "cub": "\x1b[%p1%dD"},
	},
}

func runTestCase(t *testing.T, c *testCase) {
	file, err := os.Open(c.filename)
	if err != nil {
		t.Errorf("%s: %v\n", c.filename, err)
		return
	}
	defer file.Close()

	ti, err := Read(file)
	if err != nil {
		t.Errorf("%s: %v\n", c.filename, err)
		return
	}

	names := ti.Names()
	if len(names) != len(c.names) {
		t.Errorf("%s: expected %d names, got %d.\n", c.filename, len(c.names), len(names))
	}

	for i := 0; i < len(names); i++ {
		if names[i] != c.names[i] {
			t.Errorf("%s: name mismatch. expected %s, got %s.\n", c.filename, c.names[i], names[i])
		}
	}

	for k, v := range c.boolCaps {
		v2, ok := ti.GetBoolCap(k)
		if !ok {
			t.Errorf("%s: missing bool cap '%s'\n", c.filename, k)
		} else if v != v2 {
			t.Errorf("%s: bool cap mismatch for '%s'. expected %v, got %v.\n", c.filename, k, v, v2)
		}
	}

	for k, v := range c.numCaps {
		v2, ok := ti.GetNumberCap(k)
		if !ok {
			t.Errorf("%s: missing number cap '%s'\n", c.filename, k)
		} else if v != v2 {
			t.Errorf("%s: number cap mismatch for '%s'. expected %v, got %v.\n", c.filename, k, v, v2)
		}
	}

	for k, v := range c.strCaps {
		v2, ok := ti.GetStringCap(k)
		if !ok {
			t.Errorf("%s: missing number cap '%s'\n", c.filename, k)
		} else if v != v2 {
			t.Errorf("%s: number cap mismatch for '%s'. expected %v, got %v.\n", c.filename, k, v, v2)
		}
	}
}

func TestRead(t *testing.T) {
	for i := range testCases {
		runTestCase(t, &testCases[i])
	}
}
