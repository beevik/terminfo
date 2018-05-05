package terminfo

import (
	"encoding/binary"
	"errors"
	"io"
	"strings"
)

// TermInfo contains all data describing a compiled terminfo entry.
type TermInfo struct {
	Names     []string
	BoolFlags []bool
	Numbers   []int
	Strings   []string

	numSize int
}

const (
	magic1  = 0432
	magic2  = 01036
	maxSize = 32768
)

// Errors returned by the terminfo package.
var (
	ErrInvalidFormat = errors.New("terminfo: invalid file format")
)

type header struct {
	Magic     int16
	NamesSize int16 // size of the names section in bytes
	BoolCount int16 // number of boolean values
	NumCount  int16 // number of numeric values
	StrCount  int16 // number of string values
	StrSize   int16 // size of the string table in bytes
}

// Read loads a compiled terminfo file from the reader.
func Read(r io.Reader) (*TermInfo, error) {
	var h header
	err := binary.Read(r, binary.LittleEndian, &h)
	if err != nil {
		return nil, err
	}

	ti := &TermInfo{}

	switch h.Magic {
	case magic1:
		ti.numSize = 2
	case magic2:
		ti.numSize = 4
	default:
		return nil, ErrInvalidFormat
	}

	size := binary.Size(h) + int(h.NamesSize) + int(h.BoolCount) +
		int(h.NumCount)*ti.numSize + int(h.StrCount*2) + int(h.StrSize)
	if size > maxSize {
		return nil, ErrInvalidFormat
	}

	b := make([]byte, h.NamesSize)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}
	if h.NamesSize > 0 && b[len(b)-1] != 0 {
		return nil, ErrInvalidFormat
	}
	ti.Names = strings.Split(string(b[:len(b)-1]), "|")

	ti.BoolFlags, err = readBools(r, int(h.BoolCount))
	if err != nil {
		return nil, err
	}

	if err = alignWord(r, h.NamesSize+h.BoolCount); err != nil {
		return nil, err
	}

	ti.Numbers, err = readNumbers(r, int(h.NumCount), ti.numSize)
	if err != nil {
		return nil, err
	}

	var strOffsets []int
	strOffsets, err = readNumbers(r, int(h.StrCount), 2)
	if err != nil {
		return nil, err
	}

	strTable := make([]byte, h.StrSize)
	_, err = io.ReadFull(r, strTable)
	if err != nil {
		return nil, err
	}

	ti.Strings = make([]string, h.StrCount)
	for i, o := range strOffsets {
		if o >= 0 && o < int(h.StrSize) {
			ti.Strings[i] = getString(strTable, o)
		}
	}

	return ti, nil
}

func getString(table []byte, offset int) string {
	end := offset
	for ; end < len(table); end++ {
		if table[end] == 0 {
			break
		}
	}
	return string(table[offset:end])
}

func alignWord(r io.Reader, offset int16) error {
	if (offset & 1) == 0 {
		return nil
	}

	b := make([]byte, 1)
	_, err := io.ReadFull(r, b)
	return err
}

func readBools(r io.Reader, n int) ([]bool, error) {
	b := make([]byte, n)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, err
	}

	bools := make([]bool, n)
	for i, v := range b {
		bools[i] = (v == 1)
	}

	return bools, nil
}

func readNumbers(r io.Reader, n, sz int) ([]int, error) {
	b := make([]byte, n*sz)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, err
	}

	nums := make([]int, n)
	switch sz {
	case 2:
		for i := 0; i < n; i++ {
			var v uint16
			for j := 0; j < 2; j++ {
				v |= uint16(b[0]) << uint(j*8)
				b = b[1:]
			}
			nums[i] = int(int16(v))
		}
	case 4:
		for i := 0; i < n; i++ {
			var v uint32
			for j := 0; j < 4; j++ {
				v |= uint32(b[0]) << uint(j*8)
				b = b[1:]
			}
			nums[i] = int(v)
		}
	}

	return nums, nil
}
