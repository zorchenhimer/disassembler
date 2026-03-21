package instr6502

import (
	"fmt"
	"strings"

	"git.zorchenhimer.com/Zorchenhimer/dasm/types"
)

type AsmLine struct {
	Instr string
	InlineComment string
	FullComment string
}

type DecodedInstr struct {
	Address uint
	Opcode byte
	Args   []byte
	Arg    int

	// inline parameters to a JSR
	Parameters []byte
	paramSize uint

	Instr  *Instruction

	lines []string
}

func (di *DecodedInstr) ParamSize() uint {
	return di.paramSize
}

func (di *DecodedInstr) InsertNewlineAfter() bool {
	//if len(di.Parameters) > 0 {
	if di.paramSize > 0 {
		return false
	}

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
	return uint(len(di.Args)+1)
}

func (di *DecodedInstr) render(lm types.LabelManager) {
	if len(di.lines) > 0 {
		return
	}

	buf := &strings.Builder{}
	oplen := di.Instr.Length()
	for i := uint(1); i < 3; i++ {
		if oplen > 1 {
			if lbl := lm.GetLabel(di.Address+i); lbl != nil {
				fmt.Fprintf(buf, "%s := * + %d", lbl.Name, i)
			}
		}
	}

	buf.WriteString(di.Op())
	argstr := di.ArgStr(lm)
	if argstr != "" {
		buf.WriteString(" ")
		buf.WriteString(argstr)
	}
	buf.WriteString("\n")

	if len(di.Parameters) > 0 {
		// TODO: inner labels
		params := []string{}
		for _, p := range di.Parameters {
			params = append(params, fmt.Sprintf("$%02X", p))
		}
		buf.WriteString(".byte "+ strings.Join(params, ", ")+"\n")
	}

	di.lines = strings.Split(buf.String(), "\n")
}

func (di *DecodedInstr) Asm(line int, lm types.LabelManager) (uint, string) {
	//di.render(lm)

	//offs := uint(0)
	//if line > 1 && len(di.Parameters) > 0 {
	//	offs = 3
	//}
	//if line >= len(di.lines) {
	//	return 0, ""
	//}
	//return offs, di.lines[line]

	switch line {
	case 0:
		//oplen := di.Instr.Length()
		//var b1, b2 *types.Label
		//if oplen > 1 { b1 = lm.GetLabel(di.Address+1) }
		//if oplen > 2 { b2 = lm.GetLabel(di.Address+2) }

		//var innerLabels string
		//if b1 != nil && b1.Name != "" {
		//	innerLabels = b1.Name+" := * + 1\n"
		//}
		//if b2 != nil && b1 != b2 && b2.Name != "" {
		//	innerLabels += b2.Name+" := * + 2\n"
		//}

		argstr := di.ArgStr(lm)
		if argstr != "" {
			return 0, di.Op() + " " + argstr
		}
		return 0, di.Op()
	case 1:
		if len(di.Parameters) == 0 {
			return 0, ""
		}

		inline := []string{}
		for _, b := range di.Parameters {
			inline = append(inline, fmt.Sprintf("$%02X", b))
		}
		return 3, ".byte "+strings.Join(inline, ", ")
	}

	return 0, ""
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
	switch ln {
	case 0:
		raw := di.Raw()
		rawstr := []string{}
		for _, r := range raw {
			rawstr = append(rawstr, fmt.Sprintf("%02X", r))
		}
		return fmt.Sprintf("%s", strings.Join(rawstr, " "))

	case 1:
		if len(di.Parameters) == 0 {
			return ""
		}

		inline := []string{}
		for _, b := range di.Parameters {
			inline = append(inline, fmt.Sprintf("%02X", b))
		}
		return strings.Join(inline, " ")

	default:
		return ""
	}
}

