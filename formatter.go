package dasm

import (
	"io"
	"fmt"
	"strings"
	//"strconv"

	"git.zorchenhimer.com/Zorchenhimer/dasm/types"
)

type Formatter struct {
	AsmCol     int
	CommentCol int
	FullCol    int

	CommentLevel types.CommentLevel

	lastNewline bool
	param bool

	w io.Writer
	lm types.LabelManager
}

func NewFormatter(w io.Writer, lm types.LabelManager) *Formatter {
	return &Formatter{ w: w, lm: lm}
}

func (f *Formatter) Write(address uint, line types.AsmLine) error {
	//var err error

	line.Prep(f.lm)

	// TODO: move these offset labels into line.Asm().  Print a warning
	//       if they have comments block comments; print inline comments.
	var b1, b2, b3 *types.Label
	if line.Length() == 2 {
		b2 = f.lm.GetLabel(address+1)
	} else if line.Length() == 3 {
		b2 = f.lm.GetLabel(address+1)
		b3 = f.lm.GetLabel(address+2)
	}
	b1 = f.lm.GetLabel(address)

	// If any of the found comments have a block comment, print them
	// before anything else.
	comment := false
	for _, lbl := range []*types.Label{b1, b2, b3} {
		if lbl == nil {
			continue
		}

		if lbl.CommentBlock != "" {
			comment = true
			cb := strings.Split(lbl.CommentBlock, "\n")
			if !f.lastNewline {
				fmt.Fprintln(f.w, "")
			}

			for _, c := range cb {
				c = strings.TrimSpace(c)
				if len(c) > 0 && c[0] == ';' {
					c = c[1:]
				}
				fmt.Fprintln(f.w, strings.TrimRight("; "+c, " \t"))
			}
		}
	}

	// Label that aligns with the OP code
	anon_lbl := 0
	if b1 != nil && b1.Name != "" {
		if !comment && !f.lastNewline {
			fmt.Fprintln(f.w, "")
		}
		name := b1.Name
		nl := "\n"
		if name == ":" {
			name = ""
			nl = ""
			anon_lbl = 1
		}
		fmt.Fprintf(f.w, "%s:%s", name, nl)
	}

	// Labels for the arguments (if they exist)
	if line.Length() > 1 && b2 != nil && b2.Name != "" {
		name := b2.Name
		if name == ":" {
			name = ""
		}
		fmt.Fprintf(f.w, "%s := * + 1\n", name)
	}

	if line.Length() > 2 && b3 != nil && b3 != b2 && b3.Name != "" {
		name := b3.Name
		if name == ":" {
			name = ""
		}
		fmt.Fprintf(f.w, "%s := * + 2\n", name)
	}

	for lnum := 0; lnum < line.LineCount(); lnum++ {
		offs, asm, comment := line.Asm(lnum, f.lm)
		if f.CommentLevel == types.Comment_None {
			comment = ""
		} else if comment != "" {
			comment = "; "+comment
		}
		var rawstr, fullcom string

		if f.CommentLevel == types.Comment_Full && (lnum == 0 || offs > 0) {
			rawstr = line.RawStr(lnum)
			if rawstr != "" {
				rawstr = " "+rawstr
			}
			fullcom = fmt.Sprintf("; %04X%s", address+offs, rawstr)
		}

		// Align the asm, comments, and verbose comments to
		// the configured columns.
		ln := ""
		if asm != "" {
			ln = strings.Repeat(" ", f.AsmCol-anon_lbl) + asm
		}

		if comment != "" {
			if len(ln) < f.CommentCol {
				ln += strings.Repeat(" ", f.CommentCol-len(ln))
			}
			ln += comment
		}

		if fullcom != "" {
			if len(ln) < f.FullCol {
				ln += strings.Repeat(" ", f.FullCol-len(ln))
			}
			ln += fullcom
		}
		fmt.Fprintln(f.w, ln)
	}

	f.lastNewline = line.InsertNewlineAfter() || f.param
	if line.InsertNewlineAfter() || f.param {
		fmt.Fprintln(f.w, "")
	}
	f.param = false

	if line.ParamSize() > 0 {
		f.param = true
	}

	return nil
}

