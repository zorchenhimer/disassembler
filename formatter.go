package dasm

import (
	"io"
	"fmt"
	"strings"
	//"strconv"

	"git.zorchenhimer.com/Zorchenhimer/dasm/types"
)

type Formatter struct {
	Indent int
	AsmWidth int

	CommentLevel types.CommentLevel

	w io.Writer
	lm types.LabelManager
}

func NewFormatter(w io.Writer, lm types.LabelManager) *Formatter {
	return &Formatter{ w: w, lm: lm}
}

func (f *Formatter) Write(address uint, line types.AsmLine) error {
	parts := []string{}
	if f.Indent > 0 {
		parts = append(parts, strings.Repeat(" ", f.Indent-1))
	}

	parts = append(parts, fmt.Sprintf("%-*s", f.AsmWidth-1, line.Asm(0, f.lm)))

	var err error
	lbl := f.lm.GetLabel(address)
	if lbl != nil && address == lbl.Address {
		if (lbl.Name != "" && lbl.Name != ":") || lbl.References > 0 {
			if lbl.Name == ":" {
				//lbl.Name = ""
				fmt.Printf("%#v\n", lbl)
			}
			_, err = fmt.Fprintln(f.w, lbl.Name+":")
			if err != nil {
				return err
			}
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

	_, err = fmt.Fprintln(f.w, strings.TrimRight(strings.Join(parts, " "), " "))
	if err != nil {
		return err
	}

	// TODO: figure out proper line count stuff
	if line.LineCount() > 1 {
		fmt.Fprintln(f.w, "")
	}

	return err
}
