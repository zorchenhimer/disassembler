package types

import (
	"fmt"
	"strings"
)

type DecodedData struct {
	Data    []int
	Stride  int
	IsWords bool
}

func (dd *DecodedData) LineCount() int {
	return 1
}

func (dd *DecodedData) Length() uint {
	return uint(len(dd.Data))
}

func (dd *DecodedData) Op() string {
	if dd.IsWords {
		return ".word"
	}
	return ".byte"
}

//func (dd *DecodedData) LineCount() int {
//	if len(dd.Data) <= dd.Stride {
//		return 1
//	}
//
//	lines := len(dd.Data) / dd.Stride
//	if len(dd.Data) % dd.Stride != 0 {
//		lines++
//	}
//	return lines
//}

// TODO: use stride (will need to account for multiline stuff in other areas.
func (dd *DecodedData) ArgStr(lm LabelManager) string {
	vals := []string{}
	if dd.IsWords {
		for _, dat := range dd.Data {
			lbl := lm.GetLabel(uint(dat))
			if lbl != nil {
				vals = append(vals, lbl.Name)
			} else {
				vals = append(vals, fmt.Sprintf("$%04X", dat))
			}
		}
	} else {
		for _, dat := range dd.Data {
			vals = append(vals, fmt.Sprintf("$%02X", dat))
		}
	}

	return strings.Join(vals, ", ")
}

func (dd *DecodedData) Asm(line int, lm LabelManager) string {
	return dd.Op() + " " + dd.ArgStr(lm)
}

func (dd *DecodedData) RawStr() string {
		//  11 22 33
	return "        "
}

type DecodedInstr struct {
	Address uint
	Opcode byte
	Args   []byte
	Arg    int

	Instr  *Instruction
}

func (di *DecodedInstr) LineCount() int {
	switch di.Instr.Name {
	case "RTS", "JMP":
		return 2
	default:
		return 1
	}
}

func (di *DecodedInstr) Length() uint {
	return uint(len(di.Args)+1)
}

func (di *DecodedInstr) Asm(line int, lm LabelManager) string {
	return di.Op() + " " + di.ArgStr(lm)
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
		// TODO: count anon labels between ref and addr
		if lbl.Name == "" && lbl.References > 0 {
			lbl.Name = ":"
		}
		argstr = strings.Replace(argstr, "{{arg}}", lbl.Name, 1)
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

func (di *DecodedInstr) RawStr() string {
	raw := di.Raw()
	rawstr := []string{}
	for _, r := range raw {
		rawstr = append(rawstr, fmt.Sprintf("%02X", r))
	}
	return fmt.Sprintf("%-8s", strings.Join(rawstr, " "))
}
