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
	"git.zorchenhimer.com/Zorchenhimer/dasm/instr_6502"
	//"git.zorchenhimer.com/Zorchenhimer/dasm/instr_sbx"
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

	var decoder types.Decoder

	switch cfg.Global.Architecture {
	case types.Arch_6502:
		decoder = instr6502.NewDecoder(lm, cfg.Global.AutoVars)
	case types.Arch_Full6502:
		//decoder = instr6502.NewDecoderUnofficial(lm)
		return fmt.Errorf("Full6502 is not implemented yet")
	case types.Arch_SbxScript:
		return fmt.Errorf("SBX script is not implemented yet")
	default:
		return fmt.Errorf("unknown architecture")
	}

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

		// instruction decoding
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

			if typ != types.Range_Code {
				index++
				continue
			}

			long_instr := false
			instr := decoder.TryInstr(address, raw[offset:])
			if instr == nil {
				long_instr = true
			} else {

				//instr.Instr.OpLength + instr.Instr.ArgLength
				for i := uint(0); i < instr.Length(); i++ {
					// instruction runs into defined data
					if bank.AddrType(address+i) != types.Range_Code {
						long_instr = true
					}
				}
			}

			// is instruction too long that it runs into a data range?
			if long_instr {
				// Default stride is 8.  Grab this default from somewhere?
				bank.Decoded[address] = decoder.NewData(address, raw[offset:offset+1], 8,
					types.Display_Hexadecimal, types.Range_Bytes, false)
				index++
				continue
			} else if instr.ParamSize() > 0 {
				rng := &types.Range{
					Address: address+instr.Length(),
					Size: instr.ParamSize(),
					End: address+instr.Length()+instr.ParamSize(),
					Stride: 8,
					Type: types.Range_Bytes,
					Display: types.Display_Hexadecimal,
				}
				for i := uint(0); i < instr.ParamSize(); i++ {
					bank.Ranges[i+rng.Address] = rng
				}
			}

			for i := uint(0); i < instr.Length(); i++ {
				bank.Decoded[address+i] = instr
			}

			index += instr.Length()
		}

		// TODO: Join auto-labels for each byte in word ranges.  (ie, lblA is low
		//       byte of word and lblB is high byte of word)
		/*

		Labels [
			{ address $800D; name "lbl_800D" }
		]

		; $8000
			lda L800C+0
			sta $2006
			lda lbl_800D
			sta $2006

		L800C:
		lbl_800D := * + 1
			.word $9012

		*/

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

			dd := decoder.NewData(address, raw[offset:offset+rng.Size], int(rng.Stride),
						   rng.Display, rng.Type, rng.RtsLabels)

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
		fmt.Printf("Warning: Not writing Global RAM labels anywhere!\n")
	}

	for _, bank := range cfg.Banks {
		if bank.Output == "" {
			if bank.Input == "-" {
				if len(bank.Labels) > 0 {
					fmt.Printf("Warning: Not writing RAM labels for bank %q anywhere!\n", bank.Name)
				}
			} else {
				fmt.Printf("Warning: Not writing disassembly for bank %q anywhere!\n", bank.Name)
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
		formatter.CommentLevel = cfg.Global.Comments
		formatter.AsmCol = cfg.Global.AsmColumn
		formatter.CommentCol = cfg.Global.CommentColumn
		formatter.FullCol = cfg.Global.VerboseColumn

		for addr := bank.Address; addr < bank.Address + bank.Size; {
			for _, win := range bank.Windows[addr] {
				lm.SetWindow(win.Window, win.Bank)
			}

			dec := bank.Decoded[addr]
			if dec == nil {
				panic(fmt.Sprintf("[%s] %04X+%04X (%X) how in the fuck?",
					bank.Name, bank.Address, addr, bank.Address + bank.Size))
			}

			formatter.Write(addr, dec)
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
