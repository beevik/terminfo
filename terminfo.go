package terminfo

import (
	"encoding/binary"
	"errors"
	"io"
	"strings"
)

// TermInfo contains all data describing a compiled terminfo entry.
type TermInfo struct {
	names      []string
	boolCaps   map[string]bool
	numCaps    map[string]int
	stringCaps map[string]string
	numSize    int
}

// Names returns a slice containing all the names of the TermInfo.
func (ti *TermInfo) Names() []string {
	return ti.names
}

// GetBoolCap returns a named boolean capability. You should pass the short
// variable name of the capability.
func (ti *TermInfo) GetBoolCap(name string) (v, ok bool) {
	v, ok = ti.boolCaps[name]
	return v, ok
}

// GetNumberCap returns the value of a named numeric capabilitty. You should
// pass the short name of the capability.
func (ti *TermInfo) GetNumberCap(name string) (v int, ok bool) {
	v, ok = ti.numCaps[name]
	return v, ok
}

// GetStringCap returns the value of a named string capability. You should
// pass the short name of the capability.
func (ti *TermInfo) GetStringCap(name string) (v string, ok bool) {
	v, ok = ti.stringCaps[name]
	return v, ok
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

type headerEx struct {
	BoolCount int16 // number of extended boolean values
	NumCount  int16 // number of extended numeric values
	StrCount  int16 // number of extended string values
	StrSize   int16 // size of the extended string table in bytes
	StrLimit  int16 // last offset of extended string table in bytes
}

// Read loads a compiled terminfo file from the reader.
func Read(r io.Reader) (*TermInfo, error) {
	var h header
	err := binary.Read(r, binary.LittleEndian, &h)
	if err != nil {
		return nil, err
	}

	if h.Magic != magic1 && h.Magic != magic2 {
		return nil, ErrInvalidFormat
	}

	size := h.dataSize()
	if size > maxSize {
		return nil, ErrInvalidFormat
	}

	ti := &TermInfo{}
	ti.numSize = h.numSize()

	b := make([]byte, h.NamesSize)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}
	if h.NamesSize > 0 && b[len(b)-1] != 0 {
		return nil, ErrInvalidFormat
	}
	ti.names = strings.Split(string(b[:len(b)-1]), "|")

	var bools []bool
	bools, err = readBools(r, int(h.BoolCount))
	if err != nil {
		return nil, err
	}
	ti.boolCaps = make(map[string]bool)
	for i, v := range bools {
		n := boolCapNames[i]
		ti.boolCaps[n] = v
	}

	if err = alignWord(r, h.NamesSize+h.BoolCount); err != nil {
		return nil, err
	}

	var nums []int
	nums, err = readNumbers(r, int(h.NumCount), ti.numSize)
	if err != nil {
		return nil, err
	}
	ti.numCaps = make(map[string]int)
	for i, v := range nums {
		if v != -1 {
			n := numCapNames[i]
			ti.numCaps[n] = v
		}
	}

	var strOffsets []int
	strOffsets, err = readNumbers(r, int(h.StrCount), 2)
	if err != nil {
		return nil, err
	}

	strBytes := make([]byte, h.StrSize)
	_, err = io.ReadFull(r, strBytes)
	if err != nil {
		return nil, err
	}

	strings := string(strBytes)
	ti.stringCaps = make(map[string]string)
	for i, o := range strOffsets {
		if o >= 0 && o < int(h.StrSize) {
			n := strCapNames[i]
			v := getString(strings, o)
			ti.stringCaps[n] = v
		}
	}

	err = alignWord(r, int16(size))
	if err != nil {
		return ti, nil
	}

	var h2 headerEx
	err = binary.Read(r, binary.LittleEndian, &h2)
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return ti, nil
	}
	if err != nil {
		return nil, err
	}

	var exBools []bool
	exBools, err = readBools(r, int(h2.BoolCount))
	if err != nil {
		return nil, err
	}

	err = alignWord(r, h2.BoolCount)
	if err != nil {
		return nil, err
	}

	var exNumbers []int
	exNumbers, err = readNumbers(r, int(h2.NumCount), ti.numSize)
	if err != nil {
		return nil, err
	}

	strOffsets, err = readNumbers(r, int(h2.StrCount), 2)
	if err != nil {
		return nil, err
	}

	nameCount := int(h2.BoolCount + h2.NumCount + h2.StrCount)
	var nameOffsets []int
	nameOffsets, err = readNumbers(r, nameCount, 2)
	if err != nil {
		return nil, err
	}

	strBytes = make([]byte, h2.StrLimit)
	_, err = io.ReadFull(r, strBytes)
	if err != nil {
		return nil, err
	}
	strings = string(strBytes)

	tmpStrings := make([]string, h2.StrCount)
	namesOffset := 0
	for i, o := range strOffsets {
		if o >= 0 && o < int(h2.StrLimit) {
			tmpStrings[i] = getString(strings, o)
			namesOffset = max(namesOffset, o)
		}
	}

	// Find offset where extended cap names begin.
	for ; namesOffset < len(strings); namesOffset++ {
		if strings[namesOffset] == 0 {
			namesOffset++
			break
		}
	}
	names := strings[namesOffset:]

	for _, v := range exBools {
		n := getString(names, nameOffsets[0])
		ti.boolCaps[n] = v
		nameOffsets = nameOffsets[1:]
	}

	for _, v := range exNumbers {
		n := getString(names, nameOffsets[0])
		ti.numCaps[n] = v
		nameOffsets = nameOffsets[1:]
	}

	for _, v := range tmpStrings {
		if v != "" {
			n := getString(names, nameOffsets[0])
			ti.stringCaps[n] = v
		}
		nameOffsets = nameOffsets[1:]
	}

	return ti, nil
}

func (h *header) numSize() int {
	if h.Magic == magic2 {
		return 4
	}
	return 2
}

func (h *header) dataSize() int {
	adj := 0
	if ((h.NamesSize + h.BoolCount) & 1) == 1 {
		adj = 1
	}

	return 2*6 +
		int(h.NamesSize) +
		int(h.BoolCount) +
		adj +
		int(h.NumCount)*h.numSize() +
		int(h.StrCount)*2 +
		int(h.StrSize)
}

func getString(strings string, offset int) string {
	end := offset
	for ; end < len(strings); end++ {
		if strings[end] == 0 {
			break
		}
	}
	return strings[offset:end]
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
