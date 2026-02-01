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

func (f *Formatter) Write(address uint, line types.AsmLine, lastNewline bool) error {
	var err error
	for lnum := 0; lnum < line.LineCount(); lnum++ {
		offs, asm := line.Asm(lnum, f.lm)
		lbl := f.lm.GetLabel(address+offs)

		if offs == 0 && lnum > 0 && asm == "" {
			fmt.Fprintln(f.w, "")
			continue
		}

		if lbl != nil && lbl.CommentBlock != "" && address+offs == lbl.Address {
			cb := strings.Split(lbl.CommentBlock, "\n")
			if !lastNewline {
				fmt.Fprintln(f.w, "")
			}

			for _, c := range cb {
				fmt.Fprintln(f.w, ";", strings.TrimSpace(c))
			}
		}

		parts := []string{}
		if f.Indent > 0 {
			parts = append(parts, strings.Repeat(" ", f.Indent-1))
		}

		parts = append(parts, fmt.Sprintf("%-*s", f.AsmWidth-1, asm))

		if lbl != nil && address+offs == lbl.Address {
			if (lbl.Name != "" && lbl.Name != ":") || lbl.References > 0 {
				if lbl.Name == ":" {
					_, err = fmt.Fprintln(f.w, lbl.Name)
				} else {
					_, err = fmt.Fprintln(f.w, lbl.Name+":")
				}
				if err != nil {
					return err
				}
			}
		}

		if f.CommentLevel > types.Comment_None {
			if f.CommentLevel == types.Comment_Full && asm != "" {
				parts = append(parts, ";")
				parts = append(parts, fmt.Sprintf("%04X", address+offs))
				parts = append(parts, line.RawStr())
			}

			if lbl != nil && lbl.CommentInline != "" && address+offs == lbl.Address {
				parts = append(parts, ";")
				for i, cline := range strings.Split(lbl.CommentInline, "\n") {
					if i == 0 {
						parts = append(parts, cline)
					} else {
						parts = append(parts, "\n; ", cline)
					}
				}
			}
		}

		_, err = fmt.Fprintln(f.w, strings.TrimRight(strings.Join(parts, " "), " "))
		if err != nil {
			return err
		}
	}

	if line.InsertNewlineAfter() {
		fmt.Fprintln(f.w, "")
	}

	return err
}
