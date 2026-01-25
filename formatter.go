package dasm

import (
	"io"
	"fmt"
	"strings"
	"strconv"

	"git.zorchenhimer.com/Zorchenhimer/dasm/types"
)

type Formatter struct {
	Indent int
	OpWidth int
	ArgWidth int

	CommentLevel types.CommentLevel

	w io.Writer
	lm types.LabelManager
}

func NewFormatter(w io.Writer, lm types.LabelManager) *Formatter {
	return &Formatter{ w: w, lm: lm}
}

func (f *Formatter) Write(address uint, line types.AsmLine) error {
	parts := []string{
		strings.Repeat(" ", f.Indent),
		line.Op(),
	}

	parts = append(parts, fmt.Sprintf("%"+strconv.Itoa(f.ArgWidth)+"s", line.ArgStr(f.lm)))

	lbl := f.lm.GetLabel(address)
	if lbl != nil {
		_, err := fmt.Fprintln(f.w, lbl.Name+":")
		if err != nil {
			return err
		}
	}

	if f.CommentLevel > types.Comment_None {
		parts = append(parts, ";")
		if f.CommentLevel == types.Comment_Full {
			parts = append(parts, fmt.Sprintf("%04X", address))
			parts = append(parts, line.RawStr())
		}
		if lbl != nil && lbl.Comment != "" {
			parts = append(parts, lbl.Comment)
		}
	}

	_, err := fmt.Fprintln(f.w, strings.TrimRight(strings.Join(parts, " "), " "))
	return err
}
