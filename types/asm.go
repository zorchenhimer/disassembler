package types

type AsmLine interface {
	Op() string
	ArgStr(lm LabelManager) string
	RawStr() string
	LineCount() int
	Length() uint
}
