package dasm

import (
	"fmt"
	"os"
	"path/filepath"
	//"strconv"
	//"strings"

	//"git.zorchenhimer.com/Zorchenhimer/dasm/config"
	"git.zorchenhimer.com/Zorchenhimer/dasm/instructions"
	"git.zorchenhimer.com/Zorchenhimer/dasm/types"
)

var verbose bool

func FromConfig(cfg *types.Config) error {
	abs, err := filepath.Abs(cfg.Global.Input)
	if err != nil {
		return err
	}

	raw, err := os.ReadFile(abs)
	if err != nil {
		return fmt.Errorf("Error reading input %q: %s", cfg.Global.Input, err)
	}

	lm := NewLabelManager(cfg.Global.Labels, cfg.Banks, cfg.Global.Windows)
	lm.Init()

	for _, bank := range cfg.Banks {
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

			typ := bank.Type(address)
			var decoded types.AsmLine

			if typ != types.Range_Code {
				index++
				continue
			}

			instr := instructions.TryInstr_6502(raw[offset:])
			if instr == nil {
				//bank.SetType(offs + bank.Address, types.Range_Bytes)
				dd := &types.DecodedData{
					Data: []int{int(raw[offset])},
					IsWords: false,
				}
				bank.Decoded[address] = dd
				index++
				continue
			}

			instr.Address = address //offs + bank.Address-bank.Offset
			//decoded.Instr = instr
			//decoded.Raw = instr.Raw()
			//vb("%04X %#v", offs, decoded)

			switch instr.Instr.AddrMode {
			case types.AddrMode_Accumulator,
				 types.AddrMode_Implied,
				 types.AddrMode_Immediate:
				// nope

			case types.AddrMode_Relative:
				//addr := uint((instr.Arg+1) + int(offs + bank.Address-bank.Offset))+1
				reladdr := uint(int(address) + instr.Arg + 2)
				//fmt.Printf("[%s:$%04X] rel label to $%04X (%d)\n",
				//	bank.Name, offs+bank.Address-bank.Offset, addr, instr.Arg)
				//bank.AutoLabels[addr] = types.NewLabel(addr, fmt.Sprintf("L%04X", addr))
				lm.SetLabel(types.NewLabel(reladdr, fmt.Sprintf("L%04X", reladdr)))

			default:
				lm.SetLabel(types.NewLabel(uint(instr.Arg), fmt.Sprintf("L%04X", instr.Arg)))
			}

			decoded = instr

			for i := uint(0); i < uint(instr.Instr.OpLength + instr.Instr.ArgLength); i++ {
				bank.Decoded[address+i] = decoded
			}

			index += uint(instr.Instr.OpLength + instr.Instr.ArgLength)
		}

		lm.Init()
		// ranges, specifically
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

			//fmt.Printf("processing range [%s:$%04X] %#v\n", bank.Name, address, rng)
			dd := &types.DecodedData{
				Data: []int{},
				Newline: true,
			}
			//addr := (bank.Address-bank.Offset)+offs

			if rng.Type == types.Range_Words {
				dd.IsWords = true

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

	for _, bank := range cfg.Banks {
		output, err := os.Create(bank.Output)
		if err != nil {
			return err
		}

		formatter := NewFormatter(output, lm)
		// TODO: put these in the config
		formatter.Indent = 4
		formatter.AsmWidth = 30
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

	return nil
}

func vb(format string, args ...any) {
	if !verbose {
		return
	}

	fmt.Printf(format+"\n", args...)
}
