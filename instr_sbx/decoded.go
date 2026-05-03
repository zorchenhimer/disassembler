package instrsbx

import (
	"fmt"
	"strings"

	"git.zorchenhimer.com/Zorchenhimer/dasm/types"
)

type Opcode struct {
	Address   uint
	Instr     *Instruction
	RawInline []byte
	Inline    string

	//vars []uint

	comment []string
}

func (op *Opcode) GoString() string {
	var instr string = "<nil>"
	if op.Instr != nil {
		instr = op.Instr.Name
	}
	inline := []string{}
	for _, b := range op.RawInline {
		inline = append(inline, fmt.Sprintf("$%02X", b))
	}

	return fmt.Sprintf("&instrsbx.Opcode{Address:$%04X Instr:%s RawInline:[%s] Inline:%q}",
		op.Address,
		instr,
		strings.Join(inline, " "),
		op.Inline,
	)
}

func (op *Opcode) StatName() string {
	return op.Instr.StatName()
}

func (op *Opcode) Asm(line int, lm types.LabelManager) (uint, string, string) {
	comment := ""
	if len(op.comment) >= line+1 {
		comment = op.comment[line]
	}

	switch op.Instr.Inline {
	case Inline_NullTerm:
		return 0, fmt.Sprintf("%s \"%s\"", op.Instr.Name, op.Inline), comment

	case Inline_Count, Inline_Word, Inline_Label:
		return 0, fmt.Sprintf("%s %s", op.Instr.Name, op.Inline), comment

	default:
		return 0, op.Instr.Name, comment
	}
}

func (op *Opcode) LineCount() int {
	return 1
}

func (op *Opcode) Length() uint {
	return uint(1+len(op.RawInline))
}

func (op *Opcode) Prep(lm types.LabelManager) {
	op.comment = []string{}
	lbl := lm.GetLabel(op.Address)
	if lbl != nil && lbl.CommentInline != "" {
		op.comment = strings.Split(lbl.CommentInline, "\n")
	}

	switch op.Instr.Inline {
	case Inline_NullTerm:
		chars := []string{}

		// last byte should always be the NULL byte
		for _, char := range op.RawInline[:len(op.RawInline)-1] {
			if ' ' <= char && char <= '~' {
				if char == '"' {
					chars = append(chars, `\"`)
				} else {
					chars = append(chars, string(char))
				}
			} else {
				chars = append(chars, fmt.Sprintf("\\x%02X", char))
			}
		}
		op.Inline = strings.Join(chars, "")

	case Inline_Label:
		val := uint(op.RawInline[0]) | (uint(op.RawInline[1]) << 8)
		lbl := lm.SetLabel(types.NewLabel(val, fmt.Sprintf("L%04X", val)))
		if lbl == nil {
			fmt.Printf("nil at $%04X label for $%04X\n%#v\n", op.Address, val, op)
			break
		}
		op.Inline = lbl.Name

	case Inline_Word:
		val := uint(op.RawInline[0]) | (uint(op.RawInline[1]) << 8)
		if val > 0xFF {
			op.Inline = fmt.Sprintf("$%04X", val)
		} else {
			op.Inline = fmt.Sprintf("%d", val)
		}

	case Inline_Count:
		vals := []string{}
		if len(op.RawInline) < 3 {
			fmt.Printf("missing inline bytes for Inline_Count type at $%04X\n", op.Address)
			break
		}
		vals = append(vals, fmt.Sprintf("[%d]", op.RawInline[0]))
		for i := 1; i+1 < len(op.RawInline); i += 2 {
			val := uint(op.RawInline[i]) | (uint(op.RawInline[i+1]) << 8)
			lbl := lm.SetLabel(types.NewLabel(val, fmt.Sprintf("L%04X", val)))
			if lbl == nil {
				fmt.Printf("nil at $%04X label for $%04X\n%#v\n", op.Address, val, op)
				break
			}
			vals = append(vals, lbl.Name)
		}

		op.Inline = strings.Join(vals, " ")
	}
}

func (op *Opcode) RawStr(line int) string {
	if len(op.RawInline) == 0 {
		return fmt.Sprintf("%02X", op.Instr.Opcode)
	}

	inline := []string{}
	for _, b := range op.RawInline {
		inline = append(inline, fmt.Sprintf("%02X", b))
	}

	return fmt.Sprintf("%02X %s", op.Instr.Opcode, strings.Join(inline, " "))
}

func (op *Opcode) InsertNewlineAfter() bool {
	switch op.Instr.Opcode {
	case 0x86, 0xC1, 0xAA, 0xAC, 0xFB, 0x84, // return, jump_switch, long_jump, long_return, jump_arg_a, jump_abs
		 0x81, 0x9B, 0xF2, 0xF3, 0xF4, // halts
		 0xF5, 0xF6, 0xF7, 0xF8, 0xFD:
		return true
	default:
		return false
	}
}

func (op *Opcode) ParamSize() uint {
	return 0
}

/*
*/

type StackData struct {
	Address uint
	Values  []uint

	Type types.RangeType
	Display types.RangeDisplay
	Stride int

	comment []string
}

func (dd *StackData) ArgStr(lm types.LabelManager, offset int) string {
	vals := []string{}
	for i := offset; i < len(dd.Values) && i < offset + dd.Stride; i++ {
		dat := dd.Values[i]

		if dd.Type == types.Range_Addresses {
			var lbl *types.Label
			lbl = lm.SetLabel(types.NewLabel(dat, fmt.Sprintf("L%04X", dat)))
			if lbl.Address != dat {
				offset := dat - lbl.Address
				vals = append(vals, fmt.Sprintf("%s+%d", lbl.Name, offset))
			} else {
				vals = append(vals, fmt.Sprintf("%s", lbl.Name))
			}

		} else {

			var intFmt string
			switch dd.Display {
			case types.Display_Binary:
				intFmt = "%%%08b"

			case types.Display_Decimal:
				intFmt = "%d"

			default: // types.Display_Hexadecimal
				if dd.Type == types.Range_Words {
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

func (dd *StackData) Asm(line int, lm types.LabelManager) (uint, string, string) {
	// value offset
	offs := line * dd.Stride
	argstr := dd.ArgStr(lm, offs)

	// byte offset
	if dd.Type == types.Range_Words || dd.Type == types.Range_Addresses {
		offs *= 2
	}

	return uint(offs), dd.Op() + " " + argstr, ""
}

func (dd *StackData) Op() string {
	switch dd.Type {
	case types.Range_Words:
		return ".word"
	case types.Range_Addresses:
		return ".addr"
	default:
		return ".byte"
	}
}

func (sd *StackData) InsertNewlineAfter() bool {
	return false
}

func (sd *StackData) Length() uint {
	width := 1
	if sd.Type == types.Range_Words || sd.Type == types.Range_Addresses {
		width = 2
	}
	return uint(len(sd.Values)*width)
}

// TODO: inline comments
func (sd *StackData) LineCount() int {
	if len(sd.Values) < sd.Stride {
		return 1
	}

	if sd.Stride < 1 {
		sd.Stride = 1
	}

	count := len(sd.Values) / sd.Stride
	if len(sd.Values) % sd.Stride != 0 {
		count++
	}

	return count
}

func (sd *StackData) ParamSize() uint {
	return 0
}

func (sd *StackData) Prep(lm types.LabelManager) {
	sd.comment = []string{}
	lbl := lm.GetLabel(sd.Address)
	if lbl != nil && lbl.CommentInline != "" {
		sd.comment = strings.Split(lbl.CommentInline, "\n")
	}
}

func (sd *StackData) RawStr(line int) string {
	//vals := []string{}
	//for _, v := range sd.Values {
	//	vals = append(vals, fmt.Sprintf("%02X", v))
	//}
	//return fmt.Sprintf("%v", sd.Values)
	//start := sd.Stride*ln
	//end := start + sd.Stride
	return ""
}
