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

	comment []string
}

func (op *Opcode) Asm(line int, lm types.LabelManager) (uint, string, string) {
	comment := ""
	if len(op.comment) >= line+1 {
		comment = op.comment[line]
	}

	switch op.Instr.Inline {
	case Inline_NullTerm:
		return 0, fmt.Sprintf("%s \"%s\"", op.Instr.Name, op.Inline), comment

	case Inline_CountDefault, Inline_CountNoDefault, Inline_Word:
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

	case Inline_Word:
		val := uint(op.RawInline[0]) | (uint(op.RawInline[1]) << 8)
		lbl := lm.SetLabel(types.NewLabel(val, fmt.Sprintf("L%04X", val)))
		op.Inline = lbl.Name

	case Inline_CountDefault, Inline_CountNoDefault:
		vals := []string{}
		vals = append(vals, fmt.Sprintf("%d", op.RawInline[0]))
		for i := 1; i < len(op.RawInline); i += 2 {
			val := uint(op.RawInline[i]) | (uint(op.RawInline[i+1]) << 8)
			lbl := lm.SetLabel(types.NewLabel(val, fmt.Sprintf("L%04X", val)))
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
	case 0x86, 0xC1, 0xAA, 0xAC, 0xFB, // return, jump_switch, long_jump, long_return, jump_arg_a
		 0x81, 0x9B, 0xF2, 0xF3, 0xF4, // halts
		 0xF5, 0xF6, 0xF7, 0xF8, 0xFD:
		return true
	default:
		return false
	}
}

/*
*/

func (op *Opcode) ParamSize() uint {
	return 0
}

type StackData struct {
	Address uint
	Value   int

	comment []string
}

func (sd *StackData) Asm(line int, lm types.LabelManager) (uint, string, string) {
	comment := ""
	if len(sd.comment) > line {
		comment = sd.comment[line]
	}

	if line == 0 {
		return 0, fmt.Sprintf("%d", sd.Value), comment
	}
	return 0, "", comment
}

func (sd *StackData) InsertNewlineAfter() bool {
	return false
}

func (sd *StackData) Length() uint {
	return 1
}

// TODO: inline comments
func (sd *StackData) LineCount() int {
	if len(sd.comment) > 1 {
		return len(sd.comment)
	}
	return 1
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
	return fmt.Sprintf("%02X", sd.Value)
}
