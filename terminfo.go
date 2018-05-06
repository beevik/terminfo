package terminfo

import (
	"encoding/binary"
	"errors"
	"io"
	"strings"
)

// Boolean capabilities, used to index into the TermInfo.BoolFlags slice.
const (
	BoolAutoLeftMargin = iota
	BoolAutoRightMargin
	BoolNoEscCtrlC
	BoolEOLStandoutGlitch
	BoolEatNewlineGlitch
	BoolEraseOverstrike
	BoolGenericType
	BoolHardCopy
	BoolHasMetaKey
	BoolHasStatusLine
	BoolInsertNullGlitch
	BoolMemoryAbove
	BoolMemoryBelow
	BoolMoveInsertMode
	BoolMoveStandoutMode
	BoolOverstrike
	BoolStatusLineEscOk
	BoolDestTabsMagicSMSO
	BoolTildeGlitch
	BoolTransparentUnderline
	BoolXonXoff
	BoolNeedsXonXoff
	BoolPrtrSilent
	BoolHardCursor
	BoolNonRevRmcup
	BoolNoPadChar
	BoolNonDestScrollRegion
	BoolCanChange
	BoolBackColorErase
	BoolHueLightnessSaturation
	BoolColAddrGlitch
	BoolCRCancelsMicroMode
	BoolHasPrintWheel
	BoolRowAddrGlitch
	BoolSemiAutoRightMargin
	BoolCPIChangesRes
	BoolLPIChangesRes
)

// Number capabilities, used to index into the TermInfo.Numbers slice.
const (
	NumColumns = iota
	NumInitTabs
	NumLines
	NumLinesOfMemory
	NumMagicCookieGlitch
	NumPaddingBaudRate
	NumVirtualTerminal
	NumWidthStatusLine
	NumNumLabels
	NumLabelHeight
	NumLabelWidth
	NumMaxAttributes
	NumMaximumWindows
	NumMaxColors
	NumMaxPairs
	NumNoColorVideo
	NumBufferCapacity
	NumDotVertSpacing
	NumDotHorzSpacing
	NumMaxMicroaddress
	NumMaxMicroJump
	NumMicroColSize
	NumMicroLineSize
	NumNumberOfPins
	NumOutputResChar
	NumOutputResLine
	NumOutputResHorzInch
	NumOutputResVertInch
	NumPrintRate
	NumWideCharSize
	NumButtons
	NumBitImageEntwining
	NumBitImageType
)

// String capabilities, used to index into the TermInfo.Strings slice.
const (
	StrBackTab = iota
	StrBell
	StrCarriageReturn
	StrChangeScrollRegion
	StrClearAllTabs
	StrClearScreen
	StrClrEol
	StrClrEos
	StrColumnAddress
	StrCommandCharacter
	StrCursorAddress
	StrCursorDown
	StrCursorHome
	StrCursorInvisible
	StrCursorLeft
	StrCursorMemAddress
	StrCursorNormal
	StrCursorRight
	StrCursorToLl
	StrCursorUp
	StrCursorVisible
	StrDeleteCharacter
	StrDeleteLine
	StrDisStatusLine
	StrDownHalfLine
	StrEnterAltCharsetMode
	StrEnterBlinkMode
	StrEnterBoldMode
	StrEnterCaMode
	StrEnterDeleteMode
	StrEnterDimMode
	StrEnterInsertMode
	StrEnterSecureMode
	StrEnterProtectedMode
	StrEnterReverseMode
	StrEnterStandoutMode
	StrEnterUnderlineMode
	StrEraseChars
	StrExitAltCharsetMode
	StrExitAttributeMode
	StrExitCaMode
	StrExitDeleteMode
	StrExitInsertMode
	StrExitStandoutMode
	StrExitUnderlineMode
	StrFlashScreen
	StrFormFeed
	StrFromStatusLine
	StrInit1string
	StrInit2string
	StrInit3string
	StrInitFile
	StrInsertCharacter
	StrInsertLine
	StrInsertPadding
	StrKeyBackspace
	StrKeyCatab
	StrKeyClear
	StrKeyCtab
	StrKeyDc
	StrKeyDl
	StrKeyDown
	StrKeyEic
	StrKeyEol
	StrKeyEos
	StrKeyF0
	StrKeyF1
	StrKeyF10
	StrKeyF2
	StrKeyF3
	StrKeyF4
	StrKeyF5
	StrKeyF6
	StrKeyF7
	StrKeyF8
	StrKeyF9
	StrKeyHome
	StrKeyIc
	StrKeyIl
	StrKeyLeft
	StrKeyLl
	StrKeyNpage
	StrKeyPpage
	StrKeyRight
	StrKeySf
	StrKeySr
	StrKeyStab
	StrKeyUp
	StrKeypadLocal
	StrKeypadXmit
	StrLabF0
	StrLabF1
	StrLabF10
	StrLabF2
	StrLabF3
	StrLabF4
	StrLabF5
	StrLabF6
	StrLabF7
	StrLabF8
	StrLabF9
	StrMetaOff
	StrMetaOn
	StrNewline
	StrPadChar
	StrParmDch
	StrParmDeleteLine
	StrParmDownCursor
	StrParmIch
	StrParmIndex
	StrParmInsertLine
	StrParmLeftCursor
	StrParmRightCursor
	StrParmRindex
	StrParmUpCursor
	StrPkeyKey
	StrPkeyLocal
	StrPkeyXmit
	StrPrintScreen
	StrPrtrOff
	StrPrtrOn
	StrRepeatChar
	StrReset1string
	StrReset2string
	StrReset3string
	StrResetFile
	StrRestoreCursor
	StrRowAddress
	StrSaveCursor
	StrScrollForward
	StrScrollReverse
	StrSetAttributes
	StrSetTab
	StrSetWindow
	StrTab
	StrToStatusLine
	StrUnderlineChar
	StrUpHalfLine
	StrInitProg
	StrKeyA1
	StrKeyA3
	StrKeyB2
	StrKeyC1
	StrKeyC3
	StrPrtrNon
	StrCharPadding
	StrAcsChars
	StrPlabNorm
	StrKeyBtab
	StrEnterXonMode
	StrExitXonMode
	StrEnterAmMode
	StrExitAmMode
	StrXonCharacter
	StrXoffCharacter
	StrEnaAcs
	StrLabelOn
	StrLabelOff
	StrKeyBeg
	StrKeyCancel
	StrKeyClose
	StrKeyCommand
	StrKeyCopy
	StrKeyCreate
	StrKeyEnd
	StrKeyEnter
	StrKeyExit
	StrKeyFind
	StrKeyHelp
	StrKeyMark
	StrKeyMessage
	StrKeyMove
	StrKeyNext
	StrKeyOpen
	StrKeyOptions
	StrKeyPrevious
	StrKeyPrint
	StrKeyRedo
	StrKeyReference
	StrKeyRefresh
	StrKeyReplace
	StrKeyRestart
	StrKeyResume
	StrKeySave
	StrKeySuspend
	StrKeyUndo
	StrKeySbeg
	StrKeyScancel
	StrKeyScommand
	StrKeyScopy
	StrKeyScreate
	StrKeySdc
	StrKeySdl
	StrKeySelect
	StrKeySend
	StrKeySeol
	StrKeySexit
	StrKeySfind
	StrKeyShelp
	StrKeyShome
	StrKeySic
	StrKeySleft
	StrKeySmessage
	StrKeySmove
	StrKeySnext
	StrKeySoptions
	StrKeySprevious
	StrKeySprint
	StrKeySredo
	StrKeySreplace
	StrKeySright
	StrKeySrsume
	StrKeySsave
	StrKeySsuspend
	StrKeySundo
	StrReqForInput
	StrKeyF11
	StrKeyF12
	StrKeyF13
	StrKeyF14
	StrKeyF15
	StrKeyF16
	StrKeyF17
	StrKeyF18
	StrKeyF19
	StrKeyF20
	StrKeyF21
	StrKeyF22
	StrKeyF23
	StrKeyF24
	StrKeyF25
	StrKeyF26
	StrKeyF27
	StrKeyF28
	StrKeyF29
	StrKeyF30
	StrKeyF31
	StrKeyF32
	StrKeyF33
	StrKeyF34
	StrKeyF35
	StrKeyF36
	StrKeyF37
	StrKeyF38
	StrKeyF39
	StrKeyF40
	StrKeyF41
	StrKeyF42
	StrKeyF43
	StrKeyF44
	StrKeyF45
	StrKeyF46
	StrKeyF47
	StrKeyF48
	StrKeyF49
	StrKeyF50
	StrKeyF51
	StrKeyF52
	StrKeyF53
	StrKeyF54
	StrKeyF55
	StrKeyF56
	StrKeyF57
	StrKeyF58
	StrKeyF59
	StrKeyF60
	StrKeyF61
	StrKeyF62
	StrKeyF63
	StrClrBol
	StrClearMargins
	StrSetLeftMargin
	StrSetRightMargin
	StrLabelFormat
	StrSetClock
	StrDisplayClock
	StrRemoveClock
	StrCreateWindow
	StrGotoWindow
	StrHangup
	StrDialPhone
	StrQuickDial
	StrTone
	StrPulse
	StrFlashHook
	StrFixedPause
	StrWaitTone
	StrUser0
	StrUser1
	StrUser2
	StrUser3
	StrUser4
	StrUser5
	StrUser6
	StrUser7
	StrUser8
	StrUser9
	StrOrigPair
	StrOrigColors
	StrInitializeColor
	StrInitializePair
	StrSetColorPair
	StrSetForeground
	StrSetBackground
	StrChangeCharPitch
	StrChangeLinePitch
	StrChangeResHorz
	StrChangeResVert
	StrDefineChar
	StrEnterDoublewideMode
	StrEnterDraftQuality
	StrEnterItalicsMode
	StrEnterLeftwardMode
	StrEnterMicroMode
	StrEnterNearLetterQuality
	StrEnterNormalQuality
	StrEnterShadowMode
	StrEnterSubscriptMode
	StrEnterSuperscriptMode
	StrEnterUpwardMode
	StrExitDoublewideMode
	StrExitItalicsMode
	StrExitLeftwardMode
	StrExitMicroMode
	StrExitShadowMode
	StrExitSubscriptMode
	StrExitSuperscriptMode
	StrExitUpwardMode
	StrMicroColumnAddress
	StrMicroDown
	StrMicroLeft
	StrMicroRight
	StrMicroRowAddress
	StrMicroUp
	StrOrderOfPins
	StrParmDownMicro
	StrParmLeftMicro
	StrParmRightMicro
	StrParmUpMicro
	StrSelectCharSet
	StrSetBottomMargin
	StrSetBottomMarginParm
	StrSetLeftMarginParm
	StrSetRightMarginParm
	StrSetTopMargin
	StrSetTopMarginParm
	StrStartBitImage
	StrStartCharSetDef
	StrStopBitImage
	StrStopCharSetDef
	StrSubscriptCharacters
	StrSuperscriptCharacters
	StrTheseCauseCr
	StrZeroMotion
	StrCharSetNames
	StrKeyMouse
	StrMouseInfo
	StrReqMousePos
	StrGetMouse
	StrSetAForeground
	StrSetABackground
	StrPkeyPlab
	StrDeviceType
	StrCodeSetInit
	StrSet0DesSeq
	StrSet1DesSeq
	StrSet2DesSeq
	StrSet3DesSeq
	StrSetLrMargin
	StrSetTbMargin
	StrBitImageRepeat
	StrBitImageNewline
	StrBitImageCarriageReturn
	StrColorNames
	StrDefineBitImageRegion
	StrEndBitImageRegion
	StrSetColorBand
	StrSetPageLength
	StrDisplayPcChar
	StrEnterPcCharsetMode
	StrExitPcCharsetMode
	StrEnterScancodeMode
	StrExitScancodeMode
	StrPcTermOptions
	StrScancodeEscape
	StrAltScancodeEsc
	StrEnterHorizontalHlMode
	StrEnterLeftHlMode
	StrEnterLowHlMode
	StrEnterRightHlMode
	StrEnterTopHlMode
	StrEnterVerticalHlMode
	StrSetAAttributes
	StrSetPglenInch
)

type BoolEx struct {
	Name  string
	Value bool
}

type NumberEx struct {
	Name  string
	Value int
}

type StringEx struct {
	Name  string
	Value string
}

// TermInfo contains all data describing a compiled terminfo entry.
type TermInfo struct {
	Names []string

	Bools   []bool   // Boolean capabilities
	Numbers []int    // Numeric capabilities
	Strings []string // String capabilities

	ExBools   []BoolEx   // Extended boolean capabilities
	ExNumbers []NumberEx // Extended numeric capabilities
	ExStrings []StringEx // Extended string capabilities

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
	ti.Names = strings.Split(string(b[:len(b)-1]), "|")

	ti.Bools, err = readBools(r, int(h.BoolCount))
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

	strBytes := make([]byte, h.StrSize)
	_, err = io.ReadFull(r, strBytes)
	if err != nil {
		return nil, err
	}

	strings := string(strBytes)
	ti.Strings = make([]string, h.StrCount)
	for i, o := range strOffsets {
		if o >= 0 && o < int(h.StrSize) {
			ti.Strings[i] = getString(strings, o)
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

	ti.ExStrings = make([]StringEx, h2.StrCount)
	for i, o := range strOffsets {
		if o >= 0 && o < int(h2.StrLimit) {
			ti.ExStrings[i].Value = getString(strings, o)
		}
	}

	ti.ExBools = make([]BoolEx, h2.BoolCount)
	for i, v := range exBools {
		ti.ExBools[i].Name = getString(strings, nameOffsets[0])
		ti.ExBools[i].Value = v
		nameOffsets = nameOffsets[1:]
	}

	ti.ExNumbers = make([]NumberEx, h2.NumCount)
	for i, v := range exNumbers {
		ti.ExNumbers[i].Name = getString(strings, nameOffsets[0])
		ti.ExNumbers[i].Value = v
		nameOffsets = nameOffsets[1:]
	}

	ti.ExStrings = make([]StringEx, h2.StrCount)

	_ = exBools
	_ = exNumbers
	_ = nameOffsets
	_ = strings

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
