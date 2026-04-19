package instr6502

import (
	//"fmt"
	//"strings"
	//"encoding/binary"

	//"git.zorchenhimer.com/Zorchenhimer/dasm/types"
)

type Instruction struct {
	name string

	//opLength   int
	argLength int
	//Opcode    byte
	addrMode AddrMode
}

func (i *Instruction) Length() int {
	return i.argLength + 1
}

func (i *Instruction) Name() string {
	return i.name
}

func (i *Instruction) ArgCount() int {
	if i.argLength == 0 {
		return 0
	}
	return 1
}

func (i *Instruction) ArgSize() int {
	return i.argLength
}

func (i *Instruction) StatName() string {
	return i.name +" "+ i.addrMode.String()
}

type AddrMode int

const (
	AddrMode_Absolute AddrMode = iota
	AddrMode_AbsoluteX
	AddrMode_AbsoluteY
	AddrMode_Accumulator
	AddrMode_Immediate
	AddrMode_Implied
	AddrMode_Indirect
	AddrMode_IndirectX
	AddrMode_IndirectY
	AddrMode_Relative
	AddrMode_ZeroPage
	AddrMode_ZeroPageX
	AddrMode_ZeroPageY
)

func (m AddrMode) String() string {
	switch m {
	case AddrMode_Absolute:
		return "Absolute"
	case AddrMode_AbsoluteX:
		return "Absolute X"
	case AddrMode_AbsoluteY:
		return "Absolute Y"
	case AddrMode_Accumulator:
		return "Accumulator"
	case AddrMode_Immediate:
		return "Immediate"
	case AddrMode_Implied:
		return "Implied"
	case AddrMode_Indirect:
		return "Indirect"
	case AddrMode_IndirectX:
		return "Indirect X"
	case AddrMode_IndirectY:
		return "Indirect Y"
	case AddrMode_Relative:
		return "Relative"
	case AddrMode_ZeroPage:
		return "Zero Page"
	case AddrMode_ZeroPageX:
		return "Zero Page X"
	case AddrMode_ZeroPageY:
		return "Zero Page Y"
	}

	return "Unknown"
}

var AddressModeFormats map[AddrMode]string = map[AddrMode]string{
	AddrMode_Absolute:    "{{arg}}",
	AddrMode_AbsoluteX:   "{{arg}}, X",
	AddrMode_AbsoluteY:   "{{arg}}, Y",
	AddrMode_Accumulator: "A",
	AddrMode_Immediate:   "#{{arg}}",
	AddrMode_Implied:     "",
	AddrMode_Indirect:    "({{arg}})",
	AddrMode_IndirectX:   "({{arg}}, X)",
	AddrMode_IndirectY:   "({{arg}}), Y",
	AddrMode_Relative:    "{{arg}}",
	AddrMode_ZeroPage:    "{{arg}}",
	AddrMode_ZeroPageX:   "{{arg}}, X",
	AddrMode_ZeroPageY:   "{{arg}}, Y",
}

// size in bytes
func (am AddrMode) ArgSize() int {
	switch am {
	case AddrMode_Immediate, AddrMode_Relative, AddrMode_ZeroPage,
		 AddrMode_ZeroPageX, AddrMode_ZeroPageY:
		return 1

	case AddrMode_Absolute, AddrMode_AbsoluteX, AddrMode_AbsoluteY,
		 AddrMode_Indirect, AddrMode_IndirectX, AddrMode_IndirectY:
		return 2

	default:
	//case AddrMode_Accumulator, AddrMode_Implied:
		return 0
	}
}
