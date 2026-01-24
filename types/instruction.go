package types

type Instruction struct {
	Name string

	OpLength   int
	ArgLength int
	//Opcode    byte
	AddrMode AddrMode
	AddressModeFunc AddressModeFn
}

func (i *Instruction) Length() int {
	return i.OpLength + i.ArgLength
}

// TODO: Make this an interface of some sort
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

type AddressModeFn func(arg int) string

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

