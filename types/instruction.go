package types

type Instruction interface {
	Name() string
	Length() int // in bytes; OP + args
	ArgCount() int
	ArgSize() int // in bytes
}

