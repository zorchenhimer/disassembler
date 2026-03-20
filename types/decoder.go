package types

type Decoder interface {
	TryInstr(addr uint, raw []byte) AsmLine
	NewData(addr uint, raw []byte, stride int, display RangeDisplay, rngType RangeType, rtsLabels bool) AsmLine

	SetBank(addr uint, size uint)
}
