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
	//var err error

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
			if !lastNewline {
				fmt.Fprintln(f.w, "")
			}

			for _, c := range cb {
				c = strings.TrimSpace(c)
				if len(c) > 0 && c[0] == ';' {
					c = c[1:]
				}
				fmt.Fprintln(f.w, ";", c)
			}
		}
	}

	// Label that aligns with the OP code
	if b1 != nil && b1.Name != "" {
		if !comment && !lastNewline {
			fmt.Fprintln(f.w, "")
		}
		name := b1.Name
		if name == ":" {
			name = ""
		}
		fmt.Fprintf(f.w, "%s:\n", name)
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
		offs, asm := line.Asm(lnum, f.lm)
		var rawstr, fullcom string

		if f.CommentLevel == types.Comment_Full && (lnum == 0 || offs > 0) {
			rawstr = line.RawStr(lnum)
			if rawstr != "" {
				rawstr = " "+rawstr
			}
			fullcom = fmt.Sprintf("; %04X%s", address+offs, rawstr)

			fmt.Fprintf(f.w, "%*s%-*s%s\n",
				f.Indent, "",
				f.AsmWidth, asm,
				fullcom,
			)
		} else {
			if asm == "" {
				fmt.Fprintln(f.w, "")
			} else {
				fmt.Fprintf(f.w, "%*s%s\n",
					f.Indent, "", asm,
				)
			}
		}
	}

	if line.InsertNewlineAfter() {
		fmt.Fprintln(f.w, "")
	}

	return nil
}

// FIXME: Refactor this function.  There's so much duplicated functionality
//        that breaks on edge cases, especially with inline comments and
//        intra-instruction labels.
//        Inline comments should probably be assembled inside line.Asm()
func (f *Formatter) old_Write(address uint, line types.AsmLine, lastNewline bool) error {
	var err error

	// If line as an instruction, find potential intra-instruction labels
	var b2, b3 *types.Label
	if line.Length() == 2 {
		b2 = f.lm.GetLabel(address+1)
	} else if line.Length() == 3 {
		b2 = f.lm.GetLabel(address+1)
		b3 = f.lm.GetLabel(address+2)
	}

	for lnum := 0; lnum < line.LineCount(); lnum++ {
		offs, asm := line.Asm(lnum, f.lm)
		var lbl *types.Label
		if (lnum == 0 && offs == 0) || (lnum > 0 && offs > 0) {
			lbl = f.lm.GetLabel(address+offs)
		}

		// newline for jmp and rts?
		//if offs == 0 && lnum > 0 && asm == "" {
		//	fmt.Fprintln(f.w, "")
		//	continue
		//}

		// Print comment block before anything else
		if lbl != nil && lbl.CommentBlock != "" && address+offs == lbl.Address {
			cb := strings.Split(lbl.CommentBlock, "\n")
			if !lastNewline {
				fmt.Fprintln(f.w, "")
			}

			for _, c := range cb {
				c = strings.TrimSpace(c)
				if len(c) > 0 && c[0] == ';' {
					c = c[1:]
				}
				fmt.Fprintln(f.w, ";", c)
			}
		}

		parts := []string{}
		if f.Indent > 0 {
			parts = append(parts, strings.Repeat(" ", f.Indent-1))
		}

		// Inline comments
		skipasm := false
		if f.CommentLevel > types.Comment_None {
			asmlen := len(asm)
			if lbl != nil && lbl.CommentInline != "" && address+offs == lbl.Address {
				//parts = append(parts, ";")
				if lbl.Name == ":" { // fixes "::" anon labels
					_, err = fmt.Fprintln(f.w, lbl.Name)
				} else if lbl.Name != "" {
					_, err = fmt.Fprintln(f.w, lbl.Name+":")
				}

				for i, cline := range strings.Split(lbl.CommentInline, "\n") {
					if i == 0 {
						skipasm = true
						parts = append(parts, fmt.Sprintf("%-*s", f.AsmWidth-1, asm +" ; "+ cline))
						if f.CommentLevel > types.Comment_None {
							if f.CommentLevel == types.Comment_Full && asm != "" {
								parts = append(parts, ";")
								parts = append(parts, fmt.Sprintf("%04X", address+offs))
								parts = append(parts, line.RawStr(i))
							}
						}
					} else {
						parts = append(parts, fmt.Sprintf("%-*s; %s", asmlen+1, "", cline))
					}
					_, err = fmt.Fprintln(f.w, strings.TrimRight(strings.Join(parts, " "), " "))
					if err != nil {
						return err
					}
					parts = []string{}
					if f.Indent > 0 {
						parts = append(parts, strings.Repeat(" ", f.Indent-1))
					}
				}
				continue
			}
		}

		if !skipasm {
			parts = append(parts, fmt.Sprintf("%-*s", f.AsmWidth-1, asm))
		}

		if lbl != nil && address+offs == lbl.Address {
			if (lbl.Name != "" && lbl.Name != ":") || lbl.References > 0 {
				if lbl.Name == ":" { // fixes "::" anon labels
					_, err = fmt.Fprintln(f.w, lbl.Name)
				} else {
					_, err = fmt.Fprintln(f.w, lbl.Name+":")
				}
				if err != nil {
					return err
				}
			}
		}

		// Intra-instruction labels
		if b2 != nil {
			fmt.Fprintf(f.w, "%s := * + 1\n", b2.Name)
		}
		if b3 != nil && b2 != b3 {
			fmt.Fprintf(f.w, "%s := * + 2\n", b3.Name)
		}
		b2 = nil
		b3 = nil

		if f.CommentLevel > types.Comment_None && !skipasm {
			if f.CommentLevel == types.Comment_Full && asm != "" {
				parts = append(parts, ";")
				parts = append(parts, fmt.Sprintf("%04X", address+offs))
				parts = append(parts, line.RawStr(lnum))
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
