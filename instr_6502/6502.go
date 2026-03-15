package instr6502

import (
	//"fmt"
	//"strings"
	"encoding/binary"

	"git.zorchenhimer.com/Zorchenhimer/dasm/types"
)

func Instr6502(name string, addr types.AddrMode) *types.Instruction {
	var ac int
	switch addr {
	case types.AddrMode_Accumulator, types.AddrMode_Implied:
		ac = 0
	case types.AddrMode_Immediate, types.AddrMode_Relative,
		 types.AddrMode_Indirect, types.AddrMode_IndirectX, types.AddrMode_IndirectY,
		 types.AddrMode_ZeroPage, types.AddrMode_ZeroPageX, types.AddrMode_ZeroPageY:
		ac = 1
	case types.AddrMode_Absolute, types.AddrMode_AbsoluteX, types.AddrMode_AbsoluteY:
		ac = 2
	}

	return &types.Instruction{
		Name: name,
		OpLength: 1,
		ArgLength: ac,
		AddrMode: addr,
	}
}

//var AddressModeFuncs map[AddrMode]AddressModeFn = map[AddrMode]AddressModeFn{
//	AddrMode_Absolute:    func(arg int) string { return fmt.Sprintf("$%04X",      arg) },
//	AddrMode_AbsoluteX:   func(arg int) string { return fmt.Sprintf("$%04X, X",   arg) },
//	AddrMode_AbsoluteY:   func(arg int) string { return fmt.Sprintf("$%04X, Y",   arg) },
//	AddrMode_Accumulator: func(arg int) string { return "A" },
//	AddrMode_Immediate:   func(arg int) string { return fmt.Sprintf("#%02X",      arg) },
//	AddrMode_Implied:     func(arg int) string { return "" },
//	AddrMode_Indirect:    func(arg int) string { return fmt.Sprintf("($%04X)",    arg) },
//	AddrMode_IndirectX:   func(arg int) string { return fmt.Sprintf("($%04X, X)", arg) },
//	AddrMode_IndirectY:   func(arg int) string { return fmt.Sprintf("($%04X), Y", arg) },
//	AddrMode_Relative:    func(arg int) string { return fmt.Sprintf("%d",         arg) },
//	AddrMode_ZeroPage:    func(arg int) string { return fmt.Sprintf("$%02X",      arg) },
//	AddrMode_ZeroPageX:   func(arg int) string { return fmt.Sprintf("$%02X, X",   arg) },
//	AddrMode_ZeroPageY:   func(arg int) string { return fmt.Sprintf("$%02X, Y",   arg) },
//}

var Instr_6502 map[byte]*types.Instruction = map[byte]*types.Instruction{
	OP_ADC_AB: Instr6502("ADC", types.AddrMode_Absolute),
	OP_ADC_AX: Instr6502("ADC", types.AddrMode_AbsoluteX),
	OP_ADC_AY: Instr6502("ADC", types.AddrMode_AbsoluteY),
	OP_ADC_IM: Instr6502("ADC", types.AddrMode_Immediate),
	OP_ADC_IX: Instr6502("ADC", types.AddrMode_IndirectX),
	OP_ADC_IY: Instr6502("ADC", types.AddrMode_IndirectY),
	OP_ADC_ZP: Instr6502("ADC", types.AddrMode_ZeroPage),
	OP_ADC_ZX: Instr6502("ADC", types.AddrMode_ZeroPageX),

	OP_AND_AB: Instr6502("AND", types.AddrMode_Absolute),
	OP_AND_AX: Instr6502("AND", types.AddrMode_AbsoluteX),
	OP_AND_AY: Instr6502("AND", types.AddrMode_AbsoluteY),
	OP_AND_IM: Instr6502("AND", types.AddrMode_Immediate),
	OP_AND_IX: Instr6502("AND", types.AddrMode_IndirectX),
	OP_AND_IY: Instr6502("AND", types.AddrMode_IndirectY),
	OP_AND_ZP: Instr6502("AND", types.AddrMode_ZeroPage),
	OP_AND_ZX: Instr6502("AND", types.AddrMode_ZeroPageX),

	OP_ASL_AB: Instr6502("ASL", types.AddrMode_Absolute),
	OP_ASL_AC: Instr6502("ASL", types.AddrMode_Accumulator),
	OP_ASL_AX: Instr6502("ASL", types.AddrMode_AbsoluteX),
	OP_ASL_ZP: Instr6502("ASL", types.AddrMode_ZeroPage),
	OP_ASL_ZX: Instr6502("ASL", types.AddrMode_ZeroPageX),

	OP_BCC:    Instr6502("BCC", types.AddrMode_Relative),
	OP_BCS:    Instr6502("BCS", types.AddrMode_Relative),
	OP_BEQ:    Instr6502("BEQ", types.AddrMode_Relative),
	OP_BMI:    Instr6502("BMI", types.AddrMode_Relative),
	OP_BNE:    Instr6502("BNE", types.AddrMode_Relative),
	OP_BPL:    Instr6502("BPL", types.AddrMode_Relative),
	OP_BVC:    Instr6502("BVC", types.AddrMode_Relative),
	OP_BVS:    Instr6502("BVS", types.AddrMode_Relative),

	OP_BRK:    Instr6502("BRK", types.AddrMode_Implied),

	OP_BIT_AB: Instr6502("BIT", types.AddrMode_Absolute),
	OP_BIT_ZP: Instr6502("BIT", types.AddrMode_ZeroPage),

	OP_CLC:    Instr6502("CLC", types.AddrMode_Implied),
	OP_CLD:    Instr6502("CLD", types.AddrMode_Implied),
	OP_CLI:    Instr6502("CLI", types.AddrMode_Implied),
	OP_CLV:    Instr6502("CLV", types.AddrMode_Implied),

	OP_CMP_AB: Instr6502("CMP", types.AddrMode_Absolute),
	OP_CMP_AX: Instr6502("CMP", types.AddrMode_AbsoluteX),
	OP_CMP_AY: Instr6502("CMP", types.AddrMode_AbsoluteY),
	OP_CMP_IM: Instr6502("CMP", types.AddrMode_Immediate),
	OP_CMP_IX: Instr6502("CMP", types.AddrMode_IndirectX),
	OP_CMP_IY: Instr6502("CMP", types.AddrMode_IndirectY),
	OP_CMP_ZP: Instr6502("CMP", types.AddrMode_ZeroPage),
	OP_CMP_ZX: Instr6502("CMP", types.AddrMode_ZeroPageX),

	OP_CPX_AB: Instr6502("CPX", types.AddrMode_Absolute),
	OP_CPX_IM: Instr6502("CPX", types.AddrMode_Immediate),
	OP_CPX_ZP: Instr6502("CPX", types.AddrMode_ZeroPage),

	OP_CPY_AB: Instr6502("CPY", types.AddrMode_Absolute),
	OP_CPY_IM: Instr6502("CPY", types.AddrMode_Immediate),
	OP_CPY_ZP: Instr6502("CPY", types.AddrMode_ZeroPage),

	OP_DEC_AB: Instr6502("DEC", types.AddrMode_Absolute),
	OP_DEC_AX: Instr6502("DEC", types.AddrMode_AbsoluteX),
	OP_DEC_ZP: Instr6502("DEC", types.AddrMode_ZeroPage),
	OP_DEC_ZX: Instr6502("DEC", types.AddrMode_ZeroPageX),

	OP_DEX:    Instr6502("DEX", types.AddrMode_Implied),
	OP_DEY:    Instr6502("DEY", types.AddrMode_Implied),

	OP_EOR_AB: Instr6502("EOR", types.AddrMode_Absolute),
	OP_EOR_AX: Instr6502("EOR", types.AddrMode_AbsoluteX),
	OP_EOR_AY: Instr6502("EOR", types.AddrMode_AbsoluteY),
	OP_EOR_IM: Instr6502("EOR", types.AddrMode_Immediate),
	OP_EOR_IX: Instr6502("EOR", types.AddrMode_IndirectX),
	OP_EOR_IY: Instr6502("EOR", types.AddrMode_IndirectY),
	OP_EOR_ZP: Instr6502("EOR", types.AddrMode_ZeroPage),
	OP_EOR_ZX: Instr6502("EOR", types.AddrMode_ZeroPageX),

	OP_INC_AB: Instr6502("INC", types.AddrMode_Absolute),
	OP_INC_AX: Instr6502("INC", types.AddrMode_AbsoluteX),
	OP_INC_ZP: Instr6502("INC", types.AddrMode_ZeroPage),
	OP_INC_ZX: Instr6502("INC", types.AddrMode_ZeroPageX),

	OP_INX:    Instr6502("INX", types.AddrMode_Implied),
	OP_INY:    Instr6502("INY", types.AddrMode_Implied),

	OP_JMP_AB: Instr6502("JMP", types.AddrMode_Absolute),
	OP_JMP_ID: Instr6502("JMP", types.AddrMode_Indirect),
	OP_JSR:    Instr6502("JSR", types.AddrMode_Absolute),

	OP_LDA_AB: Instr6502("LDA", types.AddrMode_Absolute),
	OP_LDA_AX: Instr6502("LDA", types.AddrMode_AbsoluteX),
	OP_LDA_AY: Instr6502("LDA", types.AddrMode_AbsoluteY),
	OP_LDA_IM: Instr6502("LDA", types.AddrMode_Immediate),
	OP_LDA_IX: Instr6502("LDA", types.AddrMode_IndirectX),
	OP_LDA_IY: Instr6502("LDA", types.AddrMode_IndirectY),
	OP_LDA_ZP: Instr6502("LDA", types.AddrMode_ZeroPage),
	OP_LDA_ZX: Instr6502("LDA", types.AddrMode_ZeroPageX),

	OP_LDX_AB: Instr6502("LDX", types.AddrMode_Absolute),
	OP_LDX_AY: Instr6502("LDX", types.AddrMode_AbsoluteY),
	OP_LDX_IM: Instr6502("LDX", types.AddrMode_Immediate),
	OP_LDX_ZP: Instr6502("LDX", types.AddrMode_ZeroPage),
	OP_LDX_ZY: Instr6502("LDX", types.AddrMode_ZeroPageY),

	OP_LDY_AB: Instr6502("LDY", types.AddrMode_Absolute),
	OP_LDY_AX: Instr6502("LDY", types.AddrMode_AbsoluteX),
	OP_LDY_IM: Instr6502("LDY", types.AddrMode_Immediate),
	OP_LDY_ZP: Instr6502("LDY", types.AddrMode_ZeroPage),
	OP_LDY_ZX: Instr6502("LDY", types.AddrMode_ZeroPageX),

	OP_LSR_AB: Instr6502("LSR", types.AddrMode_Absolute),
	OP_LSR_AC: Instr6502("LSR", types.AddrMode_Accumulator),
	OP_LSR_AX: Instr6502("LSR", types.AddrMode_AbsoluteX),
	OP_LSR_ZP: Instr6502("LSR", types.AddrMode_ZeroPage),
	OP_LSR_ZX: Instr6502("LSR", types.AddrMode_ZeroPageX),

	OP_NOP:    Instr6502("NOP", types.AddrMode_Implied),

	OP_ORA_AB: Instr6502("ORA", types.AddrMode_Absolute),
	OP_ORA_AX: Instr6502("ORA", types.AddrMode_AbsoluteX),
	OP_ORA_AY: Instr6502("ORA", types.AddrMode_AbsoluteY),
	OP_ORA_IM: Instr6502("ORA", types.AddrMode_Immediate),
	OP_ORA_IX: Instr6502("ORA", types.AddrMode_IndirectX),
	OP_ORA_IY: Instr6502("ORA", types.AddrMode_IndirectY),
	OP_ORA_ZP: Instr6502("ORA", types.AddrMode_ZeroPage),
	OP_ORA_ZX: Instr6502("ORA", types.AddrMode_ZeroPageX),

	OP_PHA:    Instr6502("PHA", types.AddrMode_Implied),
	OP_PHP:    Instr6502("PHP", types.AddrMode_Implied),
	OP_PLA:    Instr6502("PLA", types.AddrMode_Implied),
	OP_PLP:    Instr6502("PLP", types.AddrMode_Implied),

	OP_ROL_AB: Instr6502("ROL", types.AddrMode_Absolute),
	OP_ROL_AC: Instr6502("ROL", types.AddrMode_Accumulator),
	OP_ROL_AX: Instr6502("ROL", types.AddrMode_AbsoluteX),
	OP_ROL_ZP: Instr6502("ROL", types.AddrMode_ZeroPage),
	OP_ROL_ZX: Instr6502("ROL", types.AddrMode_ZeroPageX),

	OP_ROR_AB: Instr6502("ROR", types.AddrMode_Absolute),
	OP_ROR_AC: Instr6502("ROR", types.AddrMode_Accumulator),
	OP_ROR_AX: Instr6502("ROR", types.AddrMode_AbsoluteX),
	OP_ROR_ZP: Instr6502("ROR", types.AddrMode_ZeroPage),
	OP_ROR_ZX: Instr6502("ROR", types.AddrMode_ZeroPageX),

	OP_RTI:    Instr6502("RTI", types.AddrMode_Implied),
	OP_RTS:    Instr6502("RTS", types.AddrMode_Implied),

	OP_SBC_AB: Instr6502("SBC", types.AddrMode_Absolute),
	OP_SBC_AX: Instr6502("SBC", types.AddrMode_AbsoluteX),
	OP_SBC_AY: Instr6502("SBC", types.AddrMode_AbsoluteY),
	OP_SBC_IM: Instr6502("SBC", types.AddrMode_Immediate),
	OP_SBC_IX: Instr6502("SBC", types.AddrMode_IndirectX),
	OP_SBC_IY: Instr6502("SBC", types.AddrMode_IndirectY),
	OP_SBC_ZP: Instr6502("SBC", types.AddrMode_ZeroPage),
	OP_SBC_ZX: Instr6502("SBC", types.AddrMode_ZeroPageX),

	OP_SEC:    Instr6502("SEC", types.AddrMode_Implied),
	OP_SED:    Instr6502("SED", types.AddrMode_Implied),
	OP_SEI:    Instr6502("SEI", types.AddrMode_Implied),

	OP_STA_AB: Instr6502("STA", types.AddrMode_Absolute),
	OP_STA_AX: Instr6502("STA", types.AddrMode_AbsoluteX),
	OP_STA_AY: Instr6502("STA", types.AddrMode_AbsoluteY),
	OP_STA_IX: Instr6502("STA", types.AddrMode_IndirectX),
	OP_STA_IY: Instr6502("STA", types.AddrMode_IndirectY),
	OP_STA_ZP: Instr6502("STA", types.AddrMode_ZeroPage),
	OP_STA_ZX: Instr6502("STA", types.AddrMode_ZeroPageX),

	OP_STX_AB: Instr6502("STX", types.AddrMode_Absolute),
	OP_STX_ZP: Instr6502("STX", types.AddrMode_ZeroPage),
	OP_STX_ZY: Instr6502("STX", types.AddrMode_ZeroPageY),

	OP_STY_AB: Instr6502("STY", types.AddrMode_Absolute),
	OP_STY_ZP: Instr6502("STY", types.AddrMode_ZeroPage),
	OP_STY_ZX: Instr6502("STY", types.AddrMode_ZeroPageX),

	OP_TAX:    Instr6502("TAX", types.AddrMode_Implied),
	OP_TAY:    Instr6502("TAY", types.AddrMode_Implied),
	OP_TSX:    Instr6502("TSX", types.AddrMode_Implied),
	OP_TXA:    Instr6502("TXA", types.AddrMode_Implied),
	OP_TXS:    Instr6502("TXS", types.AddrMode_Implied),
	OP_TYA:    Instr6502("TYA", types.AddrMode_Implied),
}


func TryOfficial(raw []byte) *types.DecodedInstr {
	if len(raw) == 0 {
		return nil
	}

	instr, ok := Instr_6502[raw[0]]
	if !ok {
		return nil
	}

	var arg int
	var args []byte

	if instr.ArgLength == 1 {
		if len(raw) < 2 {
			return nil
		}
		args = []byte{raw[1]}

		if instr.AddrMode == types.AddrMode_Relative {
			var i8 int8
			binary.Decode(raw[1:], binary.LittleEndian, &i8)
			arg = int(i8)
		} else {
			var u8 uint8
			binary.Decode(raw[1:], binary.LittleEndian, &u8)
			arg = int(u8)
		}

	} else if instr.ArgLength == 2 {
		if len(raw) < 3 {
			return nil
		}
		args = []byte{raw[1], raw[2]}

		var u16 uint16
		binary.Decode(raw[1:], binary.LittleEndian, &u16)
		arg = int(u16)
		//fmt.Printf("} u16:%X arg:%X args:[%X %X]\n", u16, arg, args[0], args[1])
	}

	return &types.DecodedInstr{
		Opcode: raw[0],
		Args:   args,
		Arg:    arg,
		Instr:  instr,
	}
}

//type DecodedArgs struct {
//	//Instr string
//
//	ArgName string
//	ArgRaw  []byte
//
//	Format string
//}
//
//func (da DecodedArgs) Asm() string {
//	if da.Format == "" && len(ArgRaw) != 0 {
//		panic("Forgot to specify a format")
//	}
//
//	if da.Format == "" {
//		return ""
//	}
//
//	arg := da.ArgName
//	if arg == "" {
//		var val int
//		_, err := binary.Decode(da.ArgRaw, binary.LittleEndian, &val)
//		if err != nil {
//			panic(fmt.Sprintf("somehow errored decoding argument: %#v: %w", da.ArgRaw, err))
//		}
//
//		if len(di.ArgRaw) == 1 {
//			arg = fmt.Sprintf("%02X", da.ArgRaw)
//		} else {
//			arg = fmt.Sprintf("%04X", da.ArgRaw)
//		}
//	}
//
//	return fmt.Sprintf(da.Format, arg)
//	//return fmt.Sprintf("%s %s", di.Instr, fmt)
//}
//
//type DecodedInstr struct {
//	Name string
//	Args &DecodedArgs
//}
//
//func (di DecodedInstr) Asm() string {
//	return fmt.Sprintf("%s %s", di.Name, di.Args.Asm())
//}
//
//type AddressingMode struct {
//	Decode  func(raw []byte) (*DecodedInstr, error)
//	ArgSize func() int
//}
//
//type Instruction interface {
//	Name() string
//	Length() int
//	Decode(name string, args []byte) (*DecodedInstr, error)
//}
//
//var Addr_Accumulator = AddressingMode{
//	ArgSize: func() int { return 0; }
//	Decode:  func(name string, args []byte) (*DecodedInstr, error) {
//		return &DecodedInstr{
//			Instr: name,
//			Args:  &DecodedArg{
//				ArgName: "A",
//				Format: "%s",
//			},
//		}
//	},
//}
//
//var Addr_Indirect = AddressingMode{
//	ArgSize: func() int { return 2; }
//	Decode:  func(name string, args []byte) (*DecodedInstr, error) {
//		return &DecodedInstr{
//			Instr: name,
//			Args:  &DecodedArg{
//				Format: "(%s)",
//			},
//		}
//	},
//}
//
//var Addr_Implied = AddressingMode{
//	ArgSize: func() int{ return 0; }
//	Decode:   func(name string, args []byte) (*DecodedInstr, error) {
//		return &DecodedInstr{
//			Instr: name,
//			Args:  &DecodedArg{ },
//		}
//	}
//}
//
//var Arch_6502 map[byte] *AddressingMode = {
//	OP_BRK: InstrImplied{
//		OpCode:      OP_BRK,
//		Instruction: "BRK",
//		Addressing:  Addr_Implied,
//	}
//}
//
//func Try(val byte) *DecodedInstruction {
//	instr, ok := Arch_6502[val]
//	if !ok {
//		return nil
//	}
//
//	return instr
//}
