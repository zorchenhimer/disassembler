package types

type AsmLine interface {
	//Op() string
	//ArgStr(lm LabelManager) string

	// returns address offset and output
	Asm(line int, lm LabelManager) (uint, string, string)
	RawStr(line int) string
	LineCount() int
	Length() uint // byte length of op code + args
	ParamSize() uint
	Prep(lm LabelManager)

	InsertNewlineAfter() bool
}
