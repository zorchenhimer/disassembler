package instr6502

import ()

/*
   implied and relative don't have sufixes

   _AB     absolute
   _AC     accumulator
   _AX     absolute, x
   _AY     absolute, y
   _ID     indirect (jmp only)
   _IM     immediate
   _IX     (indirect, x)
   _IY     (indirect), y
   _ZP     zero page
   _ZX     zero page, x
   _ZY     zero page, y
*/
const (
	// Debug OP Codes use unofficial NOPs
	OP_DEBUG  byte = 0x1A // Print CPU history line

	// Official OP Codes
	OP_ADC_AB byte = 0x6D //Absolute
	OP_ADC_AX byte = 0x7D //Absolute,X
	OP_ADC_AY byte = 0x79 //Absolute,Y
	OP_ADC_IM byte = 0x69 //Immediate
	OP_ADC_IX byte = 0x61 //(Indirect,X)
	OP_ADC_IY byte = 0x71 //(Indirect),Y
	OP_ADC_ZP byte = 0x65 //Zero Page
	OP_ADC_ZX byte = 0x75 //Zero Page,X

	OP_AND_AB byte = 0x2D //Absolute
	OP_AND_AX byte = 0x3D //Absolute,X
	OP_AND_AY byte = 0x39 //Absolute,Y
	OP_AND_IM byte = 0x29 //Immediate
	OP_AND_IX byte = 0x21 //(Indirect,X)
	OP_AND_IY byte = 0x31 //(Indirect),Y
	OP_AND_ZP byte = 0x25 //Zero Page
	OP_AND_ZX byte = 0x35 //Zero Page,X

	OP_ASL_AB byte = 0x0E //Absolute
	OP_ASL_AC byte = 0x0A //Accumulator
	OP_ASL_AX byte = 0x1E //Absolute,X
	OP_ASL_ZP byte = 0x06 //Zero Page
	OP_ASL_ZX byte = 0x16 //Zero Page,X

	OP_BCC    byte = 0x90 //
	OP_BCS    byte = 0xB0 //
	OP_BEQ    byte = 0xF0 //
	OP_BMI    byte = 0x30 //
	OP_BNE    byte = 0xD0 //
	OP_BPL    byte = 0x10 //
	OP_BVC    byte = 0x50 //
	OP_BVS    byte = 0x70 //

	OP_BRK    byte = 0x00 //

	OP_BIT_AB byte = 0x2C //Absolute
	OP_BIT_ZP byte = 0x24 //Zero Page

	OP_CLC    byte = 0x18 //
	OP_CLD    byte = 0xD8 //
	OP_CLI    byte = 0x58 //
	OP_CLV    byte = 0xB8 //

	OP_CMP_AB byte = 0xCD //Absolute
	OP_CMP_AX byte = 0xDD //Absolute,X
	OP_CMP_AY byte = 0xD9 //Absolute,Y
	OP_CMP_IM byte = 0xC9 //Immediate
	OP_CMP_IX byte = 0xC1 //(Indirect,X)
	OP_CMP_IY byte = 0xD1 //(Indirect),Y
	OP_CMP_ZP byte = 0xC5 //Zero Page
	OP_CMP_ZX byte = 0xD5 //Zero Page,X

	OP_CPX_AB byte = 0xEC //Absolute
	OP_CPX_IM byte = 0xE0 //Immediate
	OP_CPX_ZP byte = 0xE4 //Zero Page

	OP_CPY_AB byte = 0xCC //Absolute
	OP_CPY_IM byte = 0xC0 //Immediate
	OP_CPY_ZP byte = 0xC4 //Zero Page

	OP_DEC_AB byte = 0xCE //Absolute
	OP_DEC_AX byte = 0xDE //Absolute,X
	OP_DEC_ZP byte = 0xC6 //Zero Page
	OP_DEC_ZX byte = 0xD6 //Zero Page,X

	OP_DEX    byte = 0xCA //
	OP_DEY    byte = 0x88 //

	OP_EOR_AB byte = 0x4D //Absolute
	OP_EOR_AX byte = 0x5D //Absolute,X
	OP_EOR_AY byte = 0x59 //Absolute,Y
	OP_EOR_IM byte = 0x49 //Immediate
	OP_EOR_IX byte = 0x41 //(Indirect,X)
	OP_EOR_IY byte = 0x51 //(Indirect),Y
	OP_EOR_ZP byte = 0x45 //Zero Page
	OP_EOR_ZX byte = 0x55 //Zero Page,X

	OP_INC_AB byte = 0xEE //Absolute
	OP_INC_AX byte = 0xFE //Absolute,X
	OP_INC_ZP byte = 0xE6 //Zero Page
	OP_INC_ZX byte = 0xF6 //Zero Page,X

	OP_INX    byte = 0xE8 //
	OP_INY    byte = 0xC8 //

	OP_JMP_AB byte = 0x4C //Absolute
	OP_JMP_ID byte = 0x6C //Indirect
	OP_JSR    byte = 0x20 //

	OP_LDA_AB byte = 0xAD //Absolute
	OP_LDA_AX byte = 0xBD //Absolute,X
	OP_LDA_AY byte = 0xB9 //Absolute,Y
	OP_LDA_IM byte = 0xA9 //Immediate
	OP_LDA_IX byte = 0xA1 //(Indirect,X)
	OP_LDA_IY byte = 0xB1 //(Indirect),Y
	OP_LDA_ZP byte = 0xA5 //Zero Page
	OP_LDA_ZX byte = 0xB5 //Zero Page,X

	OP_LDX_AB byte = 0xAE //Absolute
	OP_LDX_AY byte = 0xBE //Absolute,Y
	OP_LDX_IM byte = 0xA2 //Immediate
	OP_LDX_ZP byte = 0xA6 //Zero Page
	OP_LDX_ZY byte = 0xB6 //Zero Page,Y

	OP_LDY_AB byte = 0xAC //Absolute
	OP_LDY_AX byte = 0xBC //Absolute,X
	OP_LDY_IM byte = 0xA0 //Immediate
	OP_LDY_ZP byte = 0xA4 //Zero Page
	OP_LDY_ZX byte = 0xB4 //Zero Page,X

	OP_LSR_AB byte = 0x4E //Absolute
	OP_LSR_AC byte = 0x4A //Accumulator
	OP_LSR_AX byte = 0x5E //Absolute,X
	OP_LSR_ZP byte = 0x46 //Zero Page
	OP_LSR_ZX byte = 0x56 //Zero Page,X

	OP_NOP    byte = 0xEA //

	OP_ORA_AB byte = 0x0D //Absolute
	OP_ORA_AX byte = 0x1D //Absolute,X
	OP_ORA_AY byte = 0x19 //Absolute,Y
	OP_ORA_IM byte = 0x09 //Immediate
	OP_ORA_IX byte = 0x01 //(Indirect,X)
	OP_ORA_IY byte = 0x11 //(Indirect),Y
	OP_ORA_ZP byte = 0x05 //Zero Page
	OP_ORA_ZX byte = 0x15 //Zero Page,X

	OP_PHA    byte = 0x48 //
	OP_PHP    byte = 0x08 //
	OP_PLA    byte = 0x68 //
	OP_PLP    byte = 0x28 //

	OP_ROL_AB byte = 0x2E //Absolute
	OP_ROL_AC byte = 0x2A //Accumulator
	OP_ROL_AX byte = 0x3E //Absolute,X
	OP_ROL_ZP byte = 0x26 //Zero Page
	OP_ROL_ZX byte = 0x36 //Zero Page,X

	OP_ROR_AB byte = 0x6E //Absolute
	OP_ROR_AC byte = 0x6A //Accumulator
	OP_ROR_AX byte = 0x7E //Absolute,X
	OP_ROR_ZP byte = 0x66 //Zero Page
	OP_ROR_ZX byte = 0x76 //Zero Page,X

	OP_RTI    byte = 0x40 //
	OP_RTS    byte = 0x60 //

	OP_SBC_AB byte = 0xED //Absolute
	OP_SBC_AX byte = 0xFD //Absolute,X
	OP_SBC_AY byte = 0xF9 //Absolute,Y
	OP_SBC_IM byte = 0xE9 //Immediate
	OP_SBC_IX byte = 0xE1 //(Indirect,X)
	OP_SBC_IY byte = 0xF1 //(Indirect),Y
	OP_SBC_ZP byte = 0xE5 //Zero Page
	OP_SBC_ZX byte = 0xF5 //Zero Page,X

	OP_SEC    byte = 0x38 //
	OP_SED    byte = 0xF8 //
	OP_SEI    byte = 0x78 //

	OP_STA_AB byte = 0x8D //Absolute
	OP_STA_AX byte = 0x9D //Absolute,X
	OP_STA_AY byte = 0x99 //Absolute,Y
	OP_STA_IX byte = 0x81 //(Indirect,X)
	OP_STA_IY byte = 0x91 //(Indirect),Y
	OP_STA_ZP byte = 0x85 //Zero Page
	OP_STA_ZX byte = 0x95 //Zero Page,X

	OP_STX_AB byte = 0x8E //Absolute
	OP_STX_ZP byte = 0x86 //Zero Page
	OP_STX_ZY byte = 0x96 //Zero Page,Y

	OP_STY_AB byte = 0x8C //Absolute
	OP_STY_ZP byte = 0x84 //Zero Page
	OP_STY_ZX byte = 0x94 //Zero Page,X

	OP_TAX    byte = 0xAA //
	OP_TAY    byte = 0xA8 //
	OP_TSX    byte = 0xBA //
	OP_TXA    byte = 0x8A //
	OP_TXS    byte = 0x9A //
	OP_TYA    byte = 0x98 //
)

func Instr6502(name string, addr AddrMode) *Instruction {
	var ac int
	switch addr {
	case AddrMode_Accumulator, AddrMode_Implied:
		ac = 0
	case AddrMode_Immediate, AddrMode_Relative,
		 AddrMode_IndirectX, AddrMode_IndirectY,
		 AddrMode_ZeroPage, AddrMode_ZeroPageX,
		 AddrMode_ZeroPageY:
		ac = 1
	case AddrMode_Absolute, AddrMode_AbsoluteX,
	     AddrMode_AbsoluteY, AddrMode_Indirect:
		ac = 2
	}

	return &Instruction{
		name: name,
		//opLength: 1,
		argLength: ac,
		addrMode: addr,
	}
}

var Instr_6502 map[byte]*Instruction = map[byte]*Instruction{
	OP_ADC_AB: Instr6502("ADC", AddrMode_Absolute),
	OP_ADC_AX: Instr6502("ADC", AddrMode_AbsoluteX),
	OP_ADC_AY: Instr6502("ADC", AddrMode_AbsoluteY),
	OP_ADC_IM: Instr6502("ADC", AddrMode_Immediate),
	OP_ADC_IX: Instr6502("ADC", AddrMode_IndirectX),
	OP_ADC_IY: Instr6502("ADC", AddrMode_IndirectY),
	OP_ADC_ZP: Instr6502("ADC", AddrMode_ZeroPage),
	OP_ADC_ZX: Instr6502("ADC", AddrMode_ZeroPageX),

	OP_AND_AB: Instr6502("AND", AddrMode_Absolute),
	OP_AND_AX: Instr6502("AND", AddrMode_AbsoluteX),
	OP_AND_AY: Instr6502("AND", AddrMode_AbsoluteY),
	OP_AND_IM: Instr6502("AND", AddrMode_Immediate),
	OP_AND_IX: Instr6502("AND", AddrMode_IndirectX),
	OP_AND_IY: Instr6502("AND", AddrMode_IndirectY),
	OP_AND_ZP: Instr6502("AND", AddrMode_ZeroPage),
	OP_AND_ZX: Instr6502("AND", AddrMode_ZeroPageX),

	OP_ASL_AB: Instr6502("ASL", AddrMode_Absolute),
	OP_ASL_AC: Instr6502("ASL", AddrMode_Accumulator),
	OP_ASL_AX: Instr6502("ASL", AddrMode_AbsoluteX),
	OP_ASL_ZP: Instr6502("ASL", AddrMode_ZeroPage),
	OP_ASL_ZX: Instr6502("ASL", AddrMode_ZeroPageX),

	OP_BCC:    Instr6502("BCC", AddrMode_Relative),
	OP_BCS:    Instr6502("BCS", AddrMode_Relative),
	OP_BEQ:    Instr6502("BEQ", AddrMode_Relative),
	OP_BMI:    Instr6502("BMI", AddrMode_Relative),
	OP_BNE:    Instr6502("BNE", AddrMode_Relative),
	OP_BPL:    Instr6502("BPL", AddrMode_Relative),
	OP_BVC:    Instr6502("BVC", AddrMode_Relative),
	OP_BVS:    Instr6502("BVS", AddrMode_Relative),

	OP_BRK:    Instr6502("BRK", AddrMode_Implied),

	OP_BIT_AB: Instr6502("BIT", AddrMode_Absolute),
	OP_BIT_ZP: Instr6502("BIT", AddrMode_ZeroPage),

	OP_CLC:    Instr6502("CLC", AddrMode_Implied),
	OP_CLD:    Instr6502("CLD", AddrMode_Implied),
	OP_CLI:    Instr6502("CLI", AddrMode_Implied),
	OP_CLV:    Instr6502("CLV", AddrMode_Implied),

	OP_CMP_AB: Instr6502("CMP", AddrMode_Absolute),
	OP_CMP_AX: Instr6502("CMP", AddrMode_AbsoluteX),
	OP_CMP_AY: Instr6502("CMP", AddrMode_AbsoluteY),
	OP_CMP_IM: Instr6502("CMP", AddrMode_Immediate),
	OP_CMP_IX: Instr6502("CMP", AddrMode_IndirectX),
	OP_CMP_IY: Instr6502("CMP", AddrMode_IndirectY),
	OP_CMP_ZP: Instr6502("CMP", AddrMode_ZeroPage),
	OP_CMP_ZX: Instr6502("CMP", AddrMode_ZeroPageX),

	OP_CPX_AB: Instr6502("CPX", AddrMode_Absolute),
	OP_CPX_IM: Instr6502("CPX", AddrMode_Immediate),
	OP_CPX_ZP: Instr6502("CPX", AddrMode_ZeroPage),

	OP_CPY_AB: Instr6502("CPY", AddrMode_Absolute),
	OP_CPY_IM: Instr6502("CPY", AddrMode_Immediate),
	OP_CPY_ZP: Instr6502("CPY", AddrMode_ZeroPage),

	OP_DEC_AB: Instr6502("DEC", AddrMode_Absolute),
	OP_DEC_AX: Instr6502("DEC", AddrMode_AbsoluteX),
	OP_DEC_ZP: Instr6502("DEC", AddrMode_ZeroPage),
	OP_DEC_ZX: Instr6502("DEC", AddrMode_ZeroPageX),

	OP_DEX:    Instr6502("DEX", AddrMode_Implied),
	OP_DEY:    Instr6502("DEY", AddrMode_Implied),

	OP_EOR_AB: Instr6502("EOR", AddrMode_Absolute),
	OP_EOR_AX: Instr6502("EOR", AddrMode_AbsoluteX),
	OP_EOR_AY: Instr6502("EOR", AddrMode_AbsoluteY),
	OP_EOR_IM: Instr6502("EOR", AddrMode_Immediate),
	OP_EOR_IX: Instr6502("EOR", AddrMode_IndirectX),
	OP_EOR_IY: Instr6502("EOR", AddrMode_IndirectY),
	OP_EOR_ZP: Instr6502("EOR", AddrMode_ZeroPage),
	OP_EOR_ZX: Instr6502("EOR", AddrMode_ZeroPageX),

	OP_INC_AB: Instr6502("INC", AddrMode_Absolute),
	OP_INC_AX: Instr6502("INC", AddrMode_AbsoluteX),
	OP_INC_ZP: Instr6502("INC", AddrMode_ZeroPage),
	OP_INC_ZX: Instr6502("INC", AddrMode_ZeroPageX),

	OP_INX:    Instr6502("INX", AddrMode_Implied),
	OP_INY:    Instr6502("INY", AddrMode_Implied),

	OP_JMP_AB: Instr6502("JMP", AddrMode_Absolute),
	OP_JMP_ID: Instr6502("JMP", AddrMode_Indirect),
	OP_JSR:    Instr6502("JSR", AddrMode_Absolute),

	OP_LDA_AB: Instr6502("LDA", AddrMode_Absolute),
	OP_LDA_AX: Instr6502("LDA", AddrMode_AbsoluteX),
	OP_LDA_AY: Instr6502("LDA", AddrMode_AbsoluteY),
	OP_LDA_IM: Instr6502("LDA", AddrMode_Immediate),
	OP_LDA_IX: Instr6502("LDA", AddrMode_IndirectX),
	OP_LDA_IY: Instr6502("LDA", AddrMode_IndirectY),
	OP_LDA_ZP: Instr6502("LDA", AddrMode_ZeroPage),
	OP_LDA_ZX: Instr6502("LDA", AddrMode_ZeroPageX),

	OP_LDX_AB: Instr6502("LDX", AddrMode_Absolute),
	OP_LDX_AY: Instr6502("LDX", AddrMode_AbsoluteY),
	OP_LDX_IM: Instr6502("LDX", AddrMode_Immediate),
	OP_LDX_ZP: Instr6502("LDX", AddrMode_ZeroPage),
	OP_LDX_ZY: Instr6502("LDX", AddrMode_ZeroPageY),

	OP_LDY_AB: Instr6502("LDY", AddrMode_Absolute),
	OP_LDY_AX: Instr6502("LDY", AddrMode_AbsoluteX),
	OP_LDY_IM: Instr6502("LDY", AddrMode_Immediate),
	OP_LDY_ZP: Instr6502("LDY", AddrMode_ZeroPage),
	OP_LDY_ZX: Instr6502("LDY", AddrMode_ZeroPageX),

	OP_LSR_AB: Instr6502("LSR", AddrMode_Absolute),
	OP_LSR_AC: Instr6502("LSR", AddrMode_Accumulator),
	OP_LSR_AX: Instr6502("LSR", AddrMode_AbsoluteX),
	OP_LSR_ZP: Instr6502("LSR", AddrMode_ZeroPage),
	OP_LSR_ZX: Instr6502("LSR", AddrMode_ZeroPageX),

	OP_NOP:    Instr6502("NOP", AddrMode_Implied),

	OP_ORA_AB: Instr6502("ORA", AddrMode_Absolute),
	OP_ORA_AX: Instr6502("ORA", AddrMode_AbsoluteX),
	OP_ORA_AY: Instr6502("ORA", AddrMode_AbsoluteY),
	OP_ORA_IM: Instr6502("ORA", AddrMode_Immediate),
	OP_ORA_IX: Instr6502("ORA", AddrMode_IndirectX),
	OP_ORA_IY: Instr6502("ORA", AddrMode_IndirectY),
	OP_ORA_ZP: Instr6502("ORA", AddrMode_ZeroPage),
	OP_ORA_ZX: Instr6502("ORA", AddrMode_ZeroPageX),

	OP_PHA:    Instr6502("PHA", AddrMode_Implied),
	OP_PHP:    Instr6502("PHP", AddrMode_Implied),
	OP_PLA:    Instr6502("PLA", AddrMode_Implied),
	OP_PLP:    Instr6502("PLP", AddrMode_Implied),

	OP_ROL_AB: Instr6502("ROL", AddrMode_Absolute),
	OP_ROL_AC: Instr6502("ROL", AddrMode_Accumulator),
	OP_ROL_AX: Instr6502("ROL", AddrMode_AbsoluteX),
	OP_ROL_ZP: Instr6502("ROL", AddrMode_ZeroPage),
	OP_ROL_ZX: Instr6502("ROL", AddrMode_ZeroPageX),

	OP_ROR_AB: Instr6502("ROR", AddrMode_Absolute),
	OP_ROR_AC: Instr6502("ROR", AddrMode_Accumulator),
	OP_ROR_AX: Instr6502("ROR", AddrMode_AbsoluteX),
	OP_ROR_ZP: Instr6502("ROR", AddrMode_ZeroPage),
	OP_ROR_ZX: Instr6502("ROR", AddrMode_ZeroPageX),

	OP_RTI:    Instr6502("RTI", AddrMode_Implied),
	OP_RTS:    Instr6502("RTS", AddrMode_Implied),

	OP_SBC_AB: Instr6502("SBC", AddrMode_Absolute),
	OP_SBC_AX: Instr6502("SBC", AddrMode_AbsoluteX),
	OP_SBC_AY: Instr6502("SBC", AddrMode_AbsoluteY),
	OP_SBC_IM: Instr6502("SBC", AddrMode_Immediate),
	OP_SBC_IX: Instr6502("SBC", AddrMode_IndirectX),
	OP_SBC_IY: Instr6502("SBC", AddrMode_IndirectY),
	OP_SBC_ZP: Instr6502("SBC", AddrMode_ZeroPage),
	OP_SBC_ZX: Instr6502("SBC", AddrMode_ZeroPageX),

	OP_SEC:    Instr6502("SEC", AddrMode_Implied),
	OP_SED:    Instr6502("SED", AddrMode_Implied),
	OP_SEI:    Instr6502("SEI", AddrMode_Implied),

	OP_STA_AB: Instr6502("STA", AddrMode_Absolute),
	OP_STA_AX: Instr6502("STA", AddrMode_AbsoluteX),
	OP_STA_AY: Instr6502("STA", AddrMode_AbsoluteY),
	OP_STA_IX: Instr6502("STA", AddrMode_IndirectX),
	OP_STA_IY: Instr6502("STA", AddrMode_IndirectY),
	OP_STA_ZP: Instr6502("STA", AddrMode_ZeroPage),
	OP_STA_ZX: Instr6502("STA", AddrMode_ZeroPageX),

	OP_STX_AB: Instr6502("STX", AddrMode_Absolute),
	OP_STX_ZP: Instr6502("STX", AddrMode_ZeroPage),
	OP_STX_ZY: Instr6502("STX", AddrMode_ZeroPageY),

	OP_STY_AB: Instr6502("STY", AddrMode_Absolute),
	OP_STY_ZP: Instr6502("STY", AddrMode_ZeroPage),
	OP_STY_ZX: Instr6502("STY", AddrMode_ZeroPageX),

	OP_TAX:    Instr6502("TAX", AddrMode_Implied),
	OP_TAY:    Instr6502("TAY", AddrMode_Implied),
	OP_TSX:    Instr6502("TSX", AddrMode_Implied),
	OP_TXA:    Instr6502("TXA", AddrMode_Implied),
	OP_TXS:    Instr6502("TXS", AddrMode_Implied),
	OP_TYA:    Instr6502("TYA", AddrMode_Implied),
}

