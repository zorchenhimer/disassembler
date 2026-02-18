package dasm

import (
	"fmt"
	"os"
	//"path/filepath"
	//"strconv"
	"strings"
	"slices"
	"cmp"

	//"git.zorchenhimer.com/Zorchenhimer/dasm/config"
	"git.zorchenhimer.com/Zorchenhimer/dasm/instructions"
	"git.zorchenhimer.com/Zorchenhimer/dasm/types"
)

var verbose bool

func FromConfig(cfg *types.Config) error {
	inputs := cfg.GetInputs()
	raws := make(map[string][]byte)
	for _, filename := range inputs {
		// explicit no input file
		if filename == "-" {
			continue
		}

		if _, ok := raws[filename]; ok {
			continue
		}

		raw, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("Error reading input %q: %s", cfg.Global.Input, err)
		}
		raws[filename] = raw
	}

	lm := NewLabelManager(cfg.Global, cfg.Banks, cfg.Global.Windows)
	lm.Init()

	for _, bank := range cfg.Banks {
		if bank.NoDasm {
			continue
		}

		var raw []byte
		if bank.Input != "" {
			raw = raws[bank.Input]
		} else if cfg.Global.Input != "" {
			raw = raws[cfg.Global.Input]
		} else {
			fmt.Println("bank %s has no input.  skipping.", bank.Name)
			continue
		}

		if bank.Size == 0 {
			bank.Size = uint(len(raw[bank.Offset:]))
		}

		lm.Init()
		for index := uint(0); index < bank.Size ; {
			// index into raw
			offset  := index + bank.Offset
			// CPU address
			address := index + bank.Address
			if offset >= uint(len(raw)) {
				return fmt.Errorf("offset(%X;%d) past end of input(%X;%d)",
					offset, offset, len(raw), len(raw))
			}

			for _, win := range bank.Windows[address] {
				lm.SetWindow(win.Window, win.Bank)
			}

			typ := bank.AddrType(address)
			var decoded types.AsmLine

			if typ != types.Range_Code {
				index++
				continue
			}

			long_instr := false
			instr := instructions.TryInstr_6502(raw[offset:])
			if instr == nil {
				long_instr = true
			} else {

				//instr.Instr.OpLength + instr.Instr.ArgLength
				for i := uint(0); i < uint(instr.Instr.OpLength + instr.Instr.ArgLength); i++ {
					// instruction runs into defined data
					if bank.AddrType(address+i) != types.Range_Code {
						long_instr = true
					}
				}
			}

			// is instruction too long that it runs into a data range?
			if long_instr {
				dd := &types.DecodedData{
					Data: []int{int(raw[offset])},
					IsWords: false,
				}
				bank.Decoded[address] = dd
				index++
				continue
			}

			length := uint(instr.Instr.OpLength + instr.Instr.ArgLength)
			instr.Address = address //offs + bank.Address-bank.Offset

			switch instr.Instr.AddrMode {
			case types.AddrMode_Accumulator,
				 types.AddrMode_Implied,
				 types.AddrMode_Immediate:
				// nope

			case types.AddrMode_Relative:
				reladdr := uint(int(address) + instr.Arg + 2)
				lm.SetLabel(types.NewLabel(reladdr, fmt.Sprintf("L%04X", reladdr)))

			default:
				var lbl *types.Label
				doLabel := false

				lblPref := "var_"
				switch instr.Instr.Name {
				case "JMP", "JSR", "BEQ", "BNE", "BPL", "BVC", "BVS", "BMI":
					lblPref = "L"
					doLabel = true
				default:
					if cfg.Global.AutoVars {
						doLabel = true
					}
				}

				if doLabel {
					lbl = lm.SetLabel(types.NewLabel(uint(instr.Arg), fmt.Sprintf(lblPref+"%04X", instr.Arg)))
				}

				if instr.Instr.Name == "JSR" && lbl != nil && lbl.ParamSize > 0 {
					if length+lbl.ParamSize+index > uint(len(raw)) {
						return fmt.Errorf("parameter for label %s out of bounds", lbl.Name)
					}

					dd := &types.DecodedData{
						Data: []int{},
						IsWords: false,
					}

					for i := uint(0); i < lbl.ParamSize; i++ {
						// +3 is for the OP + Addr
						dd.Data = append(dd.Data, int(raw[offset+i+3]))
						bank.Decoded[address+3+i] = dd
					}

					length += lbl.ParamSize
				}
			}

			decoded = instr

			for i := uint(0); i < uint(instr.Instr.OpLength + instr.Instr.ArgLength); i++ {
				bank.Decoded[address+i] = decoded
			}

			index += length
		}

		// Split ranges
		for _, rng := range bank.CfgRanges {
			if rng.Name != "" {
				continue
			}

			var lbl *types.Label
			var last *types.Label
			var crng *types.Range = rng // current range
			for i := uint(0); i < rng.Size; i++ {
				l := bank.Labels[i+rng.Address]
				split := false
				if lbl == nil && l != nil && i != 0 {
					// found first label
					split = true
					lbl = l
				} else if lbl != nil && l != lbl {
					// found new label
					split = true
					lbl = l
				}

				// Don't split if we find the same label that just split the range.
				if l == last {
					split = false
				}
				last = l

				if split {
					prng := crng
					crng = crng.Duplicate()

					prng.Size = i+rng.Address-prng.Address
					prng.End = i+rng.Address-1

					crng.Address = i+rng.Address
					crng.Size -= prng.Size
					crng.Name = ""
					crng.Comment = ""
				}

				// reassign ranges
				bank.Ranges[i+rng.Address] = crng
			}
		}

		// ranges, specifically
		lm.Init()
		for index := uint(0); index < bank.Size; {
			// index into raw
			offset  := index + bank.Offset
			// CPU address
			address := index + bank.Address

			if offset >= uint(len(raw)) {
				return fmt.Errorf("offset(%X;%d) past end of input(%X;%d)",
					offset, offset, len(raw), len(raw))
			}

			for _, win := range bank.Windows[address] {
				lm.SetWindow(win.Window, win.Bank)
			}

			rng := bank.Ranges[address] // CPU address space

			if rng == nil || rng.Type == types.Range_Code {
				index++
				continue
			}

			dd := &types.DecodedData{
				Data: []int{},
				Newline: true,
				Stride: int(rng.Stride),
			}

			if rng.Type == types.Range_Words || rng.Type == types.Range_Addresses {
				if rng.Type == types.Range_Words {
					dd.IsWords = true
				} else {
					dd.IsAddrs = true
				}

				for i := uint(0); i < rng.Size; i+=2 {
					dd.Data = append(dd.Data, int(raw[offset+i]) | (int(raw[offset+i+1]) << 8))
				}

			} else {
				for i := uint(0); i < rng.Size; i++ {
					dd.Data = append(dd.Data, int(raw[offset+i]))
				}
			}

			for i := uint(0); i < rng.Size; i++ {
				if thing, ok := bank.Decoded[i+address]; ok {
					return fmt.Errorf("Range at $%04X (starting at $%04X) in bank %s overlaps something else: %#v",
						i+address, address, bank.Name, thing)
				}
				bank.Decoded[i+address] = dd
			}

			index += rng.Size
		}

	}

	verbose = false

	if cfg.Global.Output != "" {
		output, err := os.Create(cfg.Global.Output)
		if err != nil {
			return err
		}

		list := []*types.Label{}
		longest := 0
		for addr, lbl := range lm.Global {
			if lbl.Address != addr {
				continue
			}

			if len(lbl.Name) > longest {
				longest = len(lbl.Name)
			}
			list = append(list, lbl)
		}

		slices.SortFunc(list, func(a, b *types.Label) int {
			return cmp.Compare(a.Address, b.Address)
		})

		for _, lbl := range list {
			fmt.Fprintf(output, "%-*s := $%04X\n", longest, lbl.Name, lbl.Address)
		}

		output.Close()
	} else if len(lm.Global) > 0 {
		fmt.Printf("Warning: not writing Global RAM labels anywhere!\n")
	}

	for _, bank := range cfg.Banks {
		if bank.Output == "" {
			if bank.Input == "-" {
				if len(bank.Labels) > 0 {
					fmt.Printf("Warning: not writing RAM labels for bank %q anywhere!\n", bank.Name)
				}
			} else {
				fmt.Printf("Warning: not writing disassembly for bank %q anywhere!\n", bank.Name)
			}
			continue
		}

		output, err := os.Create(bank.Output)
		if err != nil {
			return err
		}

		// output labels
		if bank.NoDasm {
			list := []*types.Label{}
			longest := 0
			for addr, lbl := range bank.Labels {
				if lbl.Address != addr {
					continue
				}

				if len(lbl.Name) > longest {
					longest = len(lbl.Name)
				}
				list = append(list, lbl)
			}

			slices.SortFunc(list, func(a, b *types.Label) int {
				return cmp.Compare(a.Address, b.Address)
			})

			for _, lbl := range list {
				fmt.Fprintf(output, "%-*s := $%04X\n", longest, lbl.Name, lbl.Address)
			}

			output.Close()
			continue
		}

		formatter := NewFormatter(output, lm)
		// TODO: put these in the config
		formatter.Indent = cfg.Global.InstrIndent
		formatter.AsmWidth = cfg.Global.AsmWidth
		formatter.CommentLevel = cfg.Global.Comments

		lastNewline := false
		for addr := bank.Address; addr < bank.Address + bank.Size; {
			for _, win := range bank.Windows[addr] {
				lm.SetWindow(win.Window, win.Bank)
			}

			dec := bank.Decoded[addr]
			if dec == nil {
				panic(fmt.Sprintf("[%s] %04X+%04X (%X) how in the fuck?",
					bank.Name, bank.Address, addr, bank.Address + bank.Size))
			}

			formatter.Write(addr, dec, lastNewline)
			lastNewline = dec.InsertNewlineAfter()
			addr += dec.Length()
		}

		output.Close()
	}

	if cfg.Global.MlbOutput != "" {
		mlb, err := os.Create(cfg.Global.MlbOutput)
		if err != nil {
			return err
		}
		defer mlb.Close()

		for addr, lbl := range lm.Global {
			if lbl.Address != addr {
				continue
			}

			t := "NesMemory"
			if addr < 0x2000 { // 0x800-0x1FFF mirrors
				t = "NesInternalRam"
			} else if addr < 0x4400 {
				t = "NesMemory"
			}
			parts := []string{t}

			if lbl.Size > 1 {
				parts = append(parts, fmt.Sprintf("%04X-%04X", lbl.Address, lbl.Address+lbl.Size-1))
			} else {
				parts = append(parts, fmt.Sprintf("%04X", lbl.Address))
			}

			parts = append(parts, lbl.Name)

			// TODO: combine these?
			if lbl.CommentBlock != "" {
				parts = append(parts, "\\n"+strings.ReplaceAll(lbl.CommentBlock, "\n", "\\n"))
			} else if lbl.CommentInline != "" {
				parts = append(parts, strings.ReplaceAll(lbl.CommentInline, "\n", "\\n"))
			}

			fmt.Fprintln(mlb, strings.Join(parts, ":"))
		}

		for _, bank := range lm.Banks {
			t := "NesPrgRom"
			if bank.NoDasm {
				t = "NesWorkRam"
			}

			for addr, lbl := range bank.Labels {
				if lbl.Address != addr {
					continue
				}

				parts := []string{t}

				if lbl.Size > 1 {
					parts = append(parts, fmt.Sprintf("%04X-%04X", bank.Offset+(lbl.Address-bank.Address), bank.Offset+(lbl.Address-bank.Address)+lbl.Size-1))
				} else {
					parts = append(parts, fmt.Sprintf("%04X", bank.Offset+(lbl.Address-bank.Address)))
				}

				parts = append(parts, lbl.Name)

				// TODO: combine these?
				if lbl.CommentBlock != "" {
					parts = append(parts, "\\n"+strings.ReplaceAll(lbl.CommentBlock, "\n", "\\n"))
				} else if lbl.CommentInline != "" {
					parts = append(parts, strings.ReplaceAll(lbl.CommentInline, "\n", "\\n"))
				}

				fmt.Fprintln(mlb, strings.Join(parts, ":"))
			}
		}
	}

	return nil
}

func vb(format string, args ...any) {
	if !verbose {
		return
	}

	fmt.Printf(format+"\n", args...)
}
