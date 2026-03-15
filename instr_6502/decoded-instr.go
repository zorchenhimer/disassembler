package instr6502

import (
	"fmt"
	"strings"

	"git.zorchenhimer.com/Zorchenhimer/dasm/types"
)

type DecodedInstr struct {
	Address uint
	Opcode byte
	Args   []byte
	Arg    int

	// inline parameters to a JSR
	Parameters []byte

	Instr  *Instruction
}

func (di *DecodedInstr) InsertNewlineAfter() bool {
	switch di.Instr.Name() {
	case "RTS", "JMP":
		return true
	default:
		return false
	}
}

func (di *DecodedInstr) LineCount() int {
	if len(di.Parameters) != 0 {
		return 2
	}
	return 1
}

func (di *DecodedInstr) Length() uint {
	return uint(len(di.Args)+len(di.Parameters)+1)
}

func (di *DecodedInstr) Asm(line int, lm types.LabelManager) (uint, string) {
	if line > 0 {
		return 0, ""
	}
	return 0, di.Op() + " " + di.ArgStr(lm)
}

func (di *DecodedInstr) Op() string {
	if di.Instr != nil {
		return di.Instr.Name()
	}

	return ""
}

func (di *DecodedInstr) ArgStr(labels types.LabelManager) string {
	var lbl *types.Label
	switch di.Instr.addrMode {
	case AddrMode_Accumulator,
		 AddrMode_Implied,
		 AddrMode_Immediate:
		// nope

	case AddrMode_Relative:
		lbl = labels.GetLabel(uint(int(di.Address)+di.Arg+2))

	default:
		lbl = labels.GetLabel(uint(di.Arg))
	}

	argstr := AddressModeFormats[di.Instr.addrMode]
	if lbl != nil {
		if lbl.Size > 1 && di.Instr.addrMode != AddrMode_Relative {
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
		if di.Instr.addrMode == AddrMode_Relative {
			// TODO: autolabel or something
			argstr = strings.Replace(argstr, "{{arg}}", fmt.Sprintf("%d", di.Arg), 1)
		} else if di.Instr.argLength == 1 {
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

