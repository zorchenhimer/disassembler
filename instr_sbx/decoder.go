package instrsbx

import (
	"fmt"

	"git.zorchenhimer.com/Zorchenhimer/dasm/types"
)

type decoder struct {
	lm types.LabelManager
	autoVars bool
}

func NewDecoder(lm types.LabelManager, autoVars bool) types.Decoder {
	return &decoder{
		lm: lm,
		autoVars: autoVars,
	}
}

// returned line is the opcode an any inline values.  stack arguments are elsewhere.
func (dec *decoder) TryInstr(addr uint, raw []byte) types.AsmLine {
	if len(raw) == 0 {
		return nil
	}

	// value high enough for opcode?
	if raw[0] & 0x80 != 0x80 {
		return &StackData{
			Address: addr,
			Value:  int(uint(raw[0])),
		}
	}

	instr, ok := Instructions[raw[0]]
	if !ok {
		panic(fmt.Sprintf("Missing SBX Script opcode %02X", raw[0]))
	}

	op := &Opcode{
		Address: addr,
		Instr: instr,
	}

	switch op.Instr.Inline {
	case Inline_NullTerm:
		found_end := false
		for i := 1; i < len(raw[1:]) && i < 33; i++ {
			op.RawInline = append(op.RawInline, raw[i])
			if raw[i] == 0x00 {
				found_end = true
				break
			}
		}
		if !found_end {
			fmt.Printf("no end found for %s at %04X",
				op.Instr.Name, op.Address)
			return nil
		}

	case Inline_Word:
		if len(raw) < 3 {
			fmt.Printf("Not enough bytes for inline word for %s at %04X",
				op.Instr.Name, op.Address)
			return nil
		}

		op.RawInline = raw[1:3]

	case Inline_CountDefault, Inline_CountNoDefault:
		// needs at least one count byte and and one inline word
		if len(raw) < 4 {
			fmt.Printf("Not enough bytes for inline word list for %s at %04X",
				op.Instr.Name, op.Address)
			return nil
		}

		count := int(raw[1])
		if len(raw[2:]) < count*2 {
			fmt.Printf("Not enough bytes for %d inline words for %s at %04X",
				count, op.Instr.Name, op.Address)
			return nil
		}
		op.RawInline = raw[2:(count*2)+2]
	}

	return op
}

func (dec *decoder) NewData(addr uint, raw []byte, stride int, display types.RangeDisplay, rngType types.RangeType, rtsLabels bool) types.AsmLine {
	return nil
}

func (dec *decoder) SetBank(addr uint, size uint) {
	// don't do anything here
	return
}
