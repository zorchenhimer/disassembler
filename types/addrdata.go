package types

import (
	"fmt"
	"strings"
)

type DecodedData struct {
	Data     []int
	Stride   int
	IsWords  bool
	IsAddrs  bool
	Newline  bool
	RtsLabel bool

	Display RangeDisplay
}

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

func (dd *DecodedData) ArgStr(lm LabelManager, offset int) string {
	vals := []string{}
	for i := offset; i < len(dd.Data) && i < offset + dd.Stride; i++ {
		dat := uint(dd.Data[i])

		if dd.IsAddrs {
			var lbl *Label
			if dd.RtsLabel {
				// RTS Trick labels (one byte before the desired destination)
				lbl = lm.SetLabel(NewLabel(dat+1, fmt.Sprintf("L%04X", dat+1)))
				if lbl.Address != dat+1 {
					// This would probably be an error with the config.  An RTS trick
					// destination into the middle of a label is most likely a mistake.
					vals = append(vals, fmt.Sprintf("%s-1%+d", lbl.Name, dat+1-lbl.Address))
				} else {
					vals = append(vals, fmt.Sprintf("%s-1", lbl.Name))
				}
			} else {
				lbl = lm.SetLabel(NewLabel(dat, fmt.Sprintf("L%04X", dat)))
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
			case Display_Binary:
				intFmt = "%%%08b"

			case Display_Decimal:
				intFmt = "%d"

			default: // Display_Hexadecimal
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

func (dd *DecodedData) Asm(line int, lm LabelManager) (uint, string) {
	// value offset
	offs := line * dd.Stride
	argstr := dd.ArgStr(lm, offs)

	// byte offset
	if dd.IsWords || dd.IsAddrs {
		offs *= 2
	}

	return uint(offs), dd.Op() + " " + argstr
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

type DecodedInstr struct {
	Address uint
	Opcode byte
	Args   []byte
	Arg    int

	Instr  *Instruction
}

func (di *DecodedInstr) InsertNewlineAfter() bool {
	switch di.Instr.Name {
	case "RTS", "JMP":
		return true
	default:
		return false
	}
}

func (di *DecodedInstr) LineCount() int {
	return 1
}

func (di *DecodedInstr) Length() uint {
	return uint(len(di.Args)+1)
}

func (di *DecodedInstr) Asm(line int, lm LabelManager) (uint, string) {
	if line > 0 {
		return 0, ""
	}
	return 0, di.Op() + " " + di.ArgStr(lm)
}

func (di *DecodedInstr) Op() string {
	if di.Instr != nil {
		return di.Instr.Name
	}

	return ""
}

func (di *DecodedInstr) ArgStr(labels LabelManager) string {
	var lbl *Label
	switch di.Instr.AddrMode {
	case AddrMode_Accumulator,
		 AddrMode_Implied,
		 AddrMode_Immediate:
		// nope

	case AddrMode_Relative:
		lbl = labels.GetLabel(uint(int(di.Address)+di.Arg+2))

	default:
		lbl = labels.GetLabel(uint(di.Arg))
	}

	argstr := AddressModeFormats[di.Instr.AddrMode]
	if lbl != nil {
		if lbl.Size > 1 && di.Instr.AddrMode != AddrMode_Relative {
			argstr = strings.Replace(argstr, "{{arg}}", fmt.Sprintf("%s+%d", lbl.Name, uint(di.Arg) - lbl.Address), 1)
		} else {

			// Count anon labels
			var lblName = lbl.Name
			if lbl.Name == ":" {
				count := 1
				dir := "+"
				if lbl.Address < di.Address {
					dir = "-"
					for i := uint(di.Address); i > lbl.Address; i-- {
						l := labels.GetLabel(i)
						if l != nil && l.Name == ":" {
							count++
						}
					}
				} else {
					for i := uint(di.Address); i < lbl.Address; i++ {
						l := labels.GetLabel(i)
						if l != nil && l.Name == ":" {
							count++
						}
					}
				}

				lblName = ":"+strings.Repeat(dir, count)
			}
			argstr = strings.Replace(argstr, "{{arg}}", lblName, 1)
		}

	} else {
		if di.Instr.AddrMode == AddrMode_Relative {
			// TODO: autolabel or something
			argstr = strings.Replace(argstr, "{{arg}}", fmt.Sprintf("%d", di.Arg), 1)
		} else if di.Instr.ArgLength == 1 {
			argstr = strings.Replace(argstr, "{{arg}}", fmt.Sprintf("$%02X", di.Arg), 1)
		} else {
			argstr = strings.Replace(argstr, "{{arg}}", fmt.Sprintf("$%04X", di.Arg), 1)
		}
	}
	return argstr
}

func (di DecodedInstr) Raw() []byte {
	return append([]byte{di.Opcode}, di.Args...)
}

func (di *DecodedInstr) RawStr(ln int) string {
	if ln > 0 {
		return ""
	}

	raw := di.Raw()
	rawstr := []string{}
	for _, r := range raw {
		rawstr = append(rawstr, fmt.Sprintf("%02X", r))
	}
	return fmt.Sprintf("%s", strings.Join(rawstr, " "))
}
