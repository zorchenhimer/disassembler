package instr6502

import (
	"encoding/binary"
	"fmt"

	"git.zorchenhimer.com/Zorchenhimer/dasm/types"
)

func NewDecoder(lm types.LabelManager, autoVars bool) types.Decoder {
	return &decoder{
		lm: lm,
		instr: Instr_6502,
		autoVars: autoVars,
	}
}

func NewDecoderUnofficial(lm types.LabelManager, autoVars bool) types.Decoder {
	panic("unofficial opcodes not implemented")

	return &decoder{
		lm: lm,
		instr: Instr_6502,
		autoVars: autoVars,
	}
}

type decoder struct {
	lm types.LabelManager
	instr map[byte]*Instruction
	autoVars bool

	bankAddr uint
	bankSize uint
}

func (dec *decoder) InstrNames() []string {
	names := []string{}
	for _, instr := range Instr_6502 {
		names = append(names, instr.StatName())
	}
	return names
}

func (dec *decoder) SetBank(addr uint, size uint) {
	dec.bankAddr = addr
	dec.bankSize = size
}

func (dec *decoder) TryInstr(addr uint, raw []byte) types.AsmLine {
	if len(raw) == 0 {
		return nil
	}

	instr, ok := dec.instr[raw[0]]
	if !ok {
		return nil
	}

	var arg int
	var args []byte

	if instr.argLength == 1 {
		if len(raw) < 2 {
			return nil
		}
		args = []byte{raw[1]}

		if instr.addrMode == AddrMode_Relative {
			var i8 int8
			binary.Decode(raw[1:], binary.LittleEndian, &i8)
			arg = int(i8)
		} else {
			var u8 uint8
			binary.Decode(raw[1:], binary.LittleEndian, &u8)
			arg = int(u8)
		}

	} else if instr.argLength == 2 {
		if len(raw) < 3 {
			return nil
		}
		args = []byte{raw[1], raw[2]}

		var u16 uint16
		binary.Decode(raw[1:], binary.LittleEndian, &u16)
		arg = int(u16)
	}

	//var params []byte
	var paramSize uint

	switch instr.addrMode {
	case AddrMode_Accumulator,
		 AddrMode_Implied,
		 AddrMode_Immediate:
		// nope

	case AddrMode_Relative:
		reladdr := uint(int(addr) + arg + 2)
		//lblRts := ""
		//if reladdr > bank.Address && reladdr < bank.Address + bank.Size -1 && raw[reladdr - bank.Address + bank.Offset] == 0x60 {
		//	lblRts = "_rts"
		//}
		//lm.SetLabel(types.NewLabel(reladdr, fmt.Sprintf("L%04X%s", reladdr, lblRts)))
		dec.lm.SetLabel(types.NewLabel(reladdr, fmt.Sprintf("L%04X", reladdr)))


	default:
		var lbl *types.Label
		doLabel := false

		lblPref := "var_"
		switch instr.Name() {
		case "JMP", "JSR", "BEQ", "BNE", "BPL", "BVC", "BVS", "BMI":
			lblPref = "L"
			doLabel = true
		default:
			if dec.autoVars {
				doLabel = true
			}
		}

		if doLabel {
			//lblRts := ""
			//if uint(instr.Arg) > bank.Address && uint(instr.Arg) < bank.Address + bank.Size -1 && raw[uint(instr.Arg) - bank.Address + bank.Offset] == 0x60 {
			//	lblRts = "_rts"
			//}
			//lbl = lm.SetLabel(types.NewLabel(uint(instr.Arg), fmt.Sprintf(lblPref+"%04X%s", instr.Arg, lblRts)))
			lbl = dec.lm.SetLabel(types.NewLabel(uint(arg), fmt.Sprintf(lblPref+"%04X", arg)))
		}

		if (instr.Name() == "JSR" || instr.Name() == "JMP") && lbl != nil && lbl.ParamSize > 0 {
			//if length+lbl.ParamSize+index > uint(len(raw)) {
			//	return fmt.Errorf("parameter for label %s out of bounds", lbl.Name)
			//}

			//dd := &types.DecodedData{
			//	Data: []int{},
			//	IsWords: false,
			//	Stride: 8,
			//}
			//dd := dec.NewData(address, raw[offset+3:offset+lbl.ParamSize+1], false, 8, types.Display_Hexadecimal, false)

			//for i := uint(0); i < lbl.ParamSize; i++ {
			//	// +3 is for the OP + Addr
			//	//dd.Data = append(dd.Data, int(raw[offset+i+3]))
			//	bank.Decoded[address+3+i] = dd
			//}

			//params = raw[3:3+lbl.ParamSize]
			paramSize = lbl.ParamSize
		}
	}

	return &DecodedInstr{
		Opcode:     raw[0],
		Args:       args,
		Arg:        arg,
		Instr:      instr,
		Address:    addr,
		//Parameters: params,
		paramSize:  paramSize,
	}
}

func (dec *decoder) NewData(addr uint, raw []byte, stride int, display types.RangeDisplay, rngType types.RangeType, rtsLabels bool) types.AsmLine {
	dd := &DecodedData{
		Data: []int{},
		Stride: stride,
		Display: display,
		RtsLabel: rtsLabels,
		Address: addr,
		Type: rngType,
	}

	if (dd.Type == types.Range_Words || dd.Type == types.Range_Addresses) && len(raw) % 2 != 0 {
		dd.Type = types.Range_Bytes
		fmt.Printf("WARNING: [$%04X] odd number of bytes for words\n", addr)
	}

	if dd.Type == types.Range_Words || dd.Type == types.Range_Addresses {
		for i := 0; i < len(raw); i += 2 {
			dd.Data = append(dd.Data, int(raw[i]) | (int(raw[i+1]) << 8))
		}
	} else {
		for i := 0; i < len(raw); i++ {
			dd.Data = append(dd.Data, int(raw[i]))
		}
	}

	if dd.Type == types.Range_Addresses {
		for _, addr := range dd.Data {
			dec.lm.SetLabel(types.NewLabel(uint(addr), fmt.Sprintf("L%04X", addr)))
		}
	}

	return dd
}
