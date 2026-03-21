package instr6502

import (
	"fmt"
	"strings"

	"git.zorchenhimer.com/Zorchenhimer/dasm/types"
)

type DecodedData struct {
	Data     []int
	Stride   int
	IsWords  bool
	IsAddrs  bool
	Newline  bool
	RtsLabel bool
	Address  uint

	Display types.RangeDisplay
	Type types.RangeType
}

func (dd *DecodedData) ParamSize() uint {
	return 0
}

func (dd *DecodedData) Prep(lm types.LabelManager) {
}

//func NewDecodedData(addr uint, raw []byte, isWords bool, stride int, display types.RangeDisplay, rtsLabels bool) types.AsmLine {
//	dd := &DecodedData{
//		Data: []int{},
//		IsWords: isWords,
//		Stride: stride,
//		Display: display,
//		RtsLabel: rtsLabels,
//		Address: addr,
//	}
//
//	if isWords && len(raw) % 2 != 0 {
//		dd.IsWords = false
//	}
//
//	if dd.IsWords {
//		for i := 0; i < len(raw); i += 2 {
//			dd.Data = append(dd.Data, int(raw[i]) | (int(raw[i+1]) << 8))
//		}
//	} else {
//		for i := 0; i < len(raw); i++ {
//			dd.Data = append(dd.Data, int(raw[i]))
//		}
//	}
//
//	return dd
//}

func (dd *DecodedData) InsertNewlineAfter() bool {
	return dd.Newline
}

func (dd *DecodedData) LineCount() int {
	if len(dd.Data) < dd.Stride {
		return 1
	}

	count := len(dd.Data) / dd.Stride
	if len(dd.Data) % dd.Stride != 0 {
		count++
	}

	return count
	//return 1
}

func (dd *DecodedData) Length() uint {
	l := uint(len(dd.Data))
	if dd.IsWords || dd.IsAddrs {
		return l*2
	}
	return l
}

func (dd *DecodedData) Op() string {
	if dd.IsWords {
		return ".word"
	} else if dd.IsAddrs {
		return ".addr"
	}
	return ".byte"
}

func (dd *DecodedData) ArgStr(lm types.LabelManager, offset int) string {
	vals := []string{}
	for i := offset; i < len(dd.Data) && i < offset + dd.Stride; i++ {
		dat := uint(dd.Data[i])

		if dd.IsAddrs {
			var lbl *types.Label
			if dd.RtsLabel {
				// RTS Trick labels (one byte before the desired destination)
				lbl = lm.SetLabel(types.NewLabel(dat+1, fmt.Sprintf("L%04X", dat+1)))
				if lbl.Address != dat+1 {
					// This would probably be an error with the config.  An RTS trick
					// destination into the middle of a label is most likely a mistake.
					vals = append(vals, fmt.Sprintf("%s-1%+d", lbl.Name, dat+1-lbl.Address))
				} else {
					vals = append(vals, fmt.Sprintf("%s-1", lbl.Name))
				}
			} else {
				lbl = lm.SetLabel(types.NewLabel(dat, fmt.Sprintf("L%04X", dat)))
				if lbl.Address != dat {
					offset := dat - lbl.Address
					vals = append(vals, fmt.Sprintf("%s+%d", lbl.Name, offset))
				} else {
					vals = append(vals, fmt.Sprintf("%s", lbl.Name))
				}
			}

		} else {

			var intFmt string
			switch dd.Display {
			case types.Display_Binary:
				intFmt = "%%%08b"

			case types.Display_Decimal:
				intFmt = "%d"

			default: // types.Display_Hexadecimal
				if dd.IsWords {
					intFmt = "$%04X"
				} else {
					intFmt = "$%02X"
				}
			}
			vals = append(vals, fmt.Sprintf(intFmt, dat))
		}
	}

	return strings.Join(vals, ", ")
}

func (dd *DecodedData) Asm(line int, lm types.LabelManager) (uint, string, string) {
	// value offset
	offs := line * dd.Stride
	argstr := dd.ArgStr(lm, offs)

	// byte offset
	if dd.IsWords || dd.IsAddrs {
		offs *= 2
	}

	return uint(offs), dd.Op() + " " + argstr, ""
}

func (dd *DecodedData) RawStr(ln int) string {
	if dd.IsAddrs {
		start := dd.Stride * ln
		end := start + dd.Stride

		if start > len(dd.Data) {
			return ""
		}

		if end > len(dd.Data) {
			end = len(dd.Data) // todo verify this
		}

		vals := []string{}
		for _, val := range dd.Data[start:end] {
			vals = append(vals, fmt.Sprintf("%04X", val))
		}
		return strings.Join(vals, " ")
	}

	return ""
}
