package types

import (
	"fmt"
	"strings"

	//"git.zorchenhimer.com/Zorchenhimer/dasm/instructions"
)

type Decoded struct {
	// nil if data
	Instr *DecodedInstr
	Val   uint

	Type RangeType
	Raw  []byte
}

func (d Decoded) GoString() string {
	if d.Instr != nil {
		return fmt.Sprintf("Decoded{Instr:%s Raw:%v}", d.Instr.Instr.Name, d.Raw)
	}

	return fmt.Sprintf("Decoded{Instr:<nil> Raw:%v", d.Raw)
}

type LabelManager interface {
	GetLabel(addr uint) *Label
	SetLabel(lbl *Label)
	GetRange(addr uint) *Range
}

func (d Decoded) Asm(addr uint, labels LabelManager) string {
	if d.Instr != nil {
		return d.Instr.Asm(addr, labels)
	}

	if d.Type == Range_Words {
		return fmt.Sprintf(".word $%04X", d.Val)
	}

	return fmt.Sprintf(".byte $%02X", d.Val)
}

type DecodedInstr struct {
	Opcode byte
	Args   []byte
	Arg    int

	Instr  *Instruction
}

func (di DecodedInstr) Raw() []byte {
	return append([]byte{di.Opcode}, di.Args...)
}

func (di DecodedInstr) Asm(addr uint, labels LabelManager) string {
	raw := di.Raw()
	rawstr := []string{}
	for _, r := range raw {
		rawstr = append(rawstr, fmt.Sprintf("%02X", r))
	}

	// TODO: proper label management
	var lbl *Label
	switch di.Instr.AddrMode {
	case AddrMode_Accumulator,
		 AddrMode_Implied,
		 AddrMode_Immediate:
		// nope

	case AddrMode_Relative:
		lbl = labels.GetLabel(uint(int(addr)+di.Arg+1))

	default:
		lbl = labels.GetLabel(uint(di.Arg))
	}

	argstr := AddressModeFormats[di.Instr.AddrMode]
	if lbl != nil {
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

	return fmt.Sprintf("%s %-10s ; %04X %s", di.Instr.Name, argstr, addr, strings.Join(rawstr, " "))
}
