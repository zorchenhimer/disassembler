package types

type AsmLine interface {
	//Op() string
	//ArgStr(lm LabelManager) string

	// returns address offset and output
	Asm(line int, lm LabelManager) (uint, string)
	RawStr(line int) string
	LineCount() int
	Length() uint

	InsertNewlineAfter() bool
}
