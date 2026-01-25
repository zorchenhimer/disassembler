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

	for _, bank := range cfg.Banks {
		//if bank.Name == "bank_01" {
		//	verbose = true
		//}

		dbg, err := os.Create("testdata/"+bank.Name+".dbg")
		if err != nil {
			return err
		}

		//for addr, rng := range bank.Ranges {
		for i := uint(bank.Address); i < bank.Address + bank.Size; i++ {
			//fmt.Fprintf(dbg, "$%04X %#v\n", addr, rng)
			fmt.Fprintf(dbg, "$%04X %s\n", i, bank.Type(i))
		}
		dbg.Close()

		//lastoffs := 0
		//for offs := bank.Offset; offs < bank.Offset + bank.Size ; {
		for index := uint(0); index < bank.Size ; {
			// index into raw
			offset  := index + bank.Offset
			// CPU address
			address := index + bank.Address
			if offset >= uint(len(raw)) {
				return fmt.Errorf("offset(%X;%d) past end of input(%X;%d)",
					offset, offset, len(raw), len(raw))
			}

			//decoded := &types.Decoded{
			//	Type: bank.Type(offs + bank.Address),
			//}
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
				//bank.AutoLabels[uint(instr.Arg)] = types.NewLabel(uint(instr.Arg), fmt.Sprintf("L%04X", instr.Arg))
				lm.SetLabel(types.NewLabel(uint(instr.Arg), fmt.Sprintf("L%04X", instr.Arg)))
			}

			decoded = instr

			for i := uint(0); i < uint(instr.Instr.OpLength + instr.Instr.ArgLength); i++ {
				//addr := bank.Address+i+(offs-bank.Offset)
				bank.Decoded[address+i] = decoded
				//vb("[%04X] %s %#v", addr, decoded.Type, decoded)
			}

			//fmt.Println(decoded.Asm(offs+bank.Address))
			index += uint(instr.Instr.OpLength + instr.Instr.ArgLength)
		}

		fmt.Println("len(bank.Ranges):", len(bank.Ranges))

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

			rng := bank.Ranges[address] // CPU address space

			if rng == nil || rng.Type == types.Range_Code {
				index++
				continue
			}

			fmt.Printf("processing range [%s:$%04X] %#v\n", bank.Name, address, rng)
			dd := &types.DecodedData{
				Data: []int{},
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
	//return nil

	// FIXME: this is borked.  Formatter isn't used and bank.Decoded probably
	//        has to change up a bit
	for _, bank := range cfg.Banks {
		output, err := os.Create(bank.Output)
		if err != nil {
			return err
		}
		formatter := NewFormatter(output, lm)

		for i := bank.Address; i < bank.Address + bank.Size; {
			dec := bank.Decoded[i]
			if dec == nil {
				panic(fmt.Sprintf("[%s] %04X+%04X (%X) how in the fuck?", bank.Name, bank.Address, i, bank.Address + bank.Size))
			}
			//lbl := bank.Labels[i]
			lbl := lm.GetLabel(i)
			if lbl != nil {
				fmt.Fprintln(output, lbl.Name+":")
			}
			//fmt.Fprintln(output, dec.Asm(i, lm))
			formatter.Write(i, dec)
			i += dec.Length()
		}

		output.Close()
	}

	//for offs := start; offs < end && offs < len(raw); {
	//	instr := instructions.TryInstr_6502(raw[offs:])
	//	if instr != nil {
	//		//fmt.Printf("%s\n", instr.Asm(offs+cfg.Banks[0].Address))
	//		fmt.Println(asm(uint(offs+cfg.Banks[0].Address), instr, state))
	//		offs += instr.Instr.Length()
	//	} else {
	//		fmt.Printf(".byte $%02X ; %04X %02X\n", raw[offs], offs+cfg.Banks[0].Address, raw[offs])
	//		//fmt.Printf(".byte $%02X\n", raw[offs])
	//		offs++
	//	}
	//	ttl--

	//	if ttl < 0 {
	//		return fmt.Errorf("TTL")
	//	}
	//}

	return nil
}

func vb(format string, args ...any) {
	if !verbose {
		return
	}

	fmt.Printf(format+"\n", args...)
}

//func Pass1(cfg *types.Config) 

//func asm(addr uint, instr *instructions.DecodedInstr, st *State) string {
//	var argstr string
//	var lbl *types.Label
//
//	if instr.Instr.AddrMode == instructions.AddrMode_Relative {
//		lbl = st.Label(addr+uint(instr.Arg+1))
//	} else {
//		lbl = st.Label(uint(addr))
//	}
//
//	if lbl != nil {
//		argstr = lbl.Name
//		if uint(lbl.Address) != addr {
//			offset := int(addr) - lbl.Address
//			argstr += "+"+strconv.Itoa(offset)
//		}
//
//	} else {
//		argSize := instr.Instr.AddrMode.ArgSize()
//
//		switch argSize {
//		case 1:
//			argstr = fmt.Sprintf("$%02X", instr.Arg)
//		case 2:
//			argstr = fmt.Sprintf("$%04X", instr.Arg)
//		}
//	}
//
//	argstr = strings.Replace(
//		instructions.AddressModeFormats[instr.Instr.AddrMode],
//		"{{arg}}", argstr, 1)
//
//	raw := append([]byte{instr.Opcode}, instr.Args...)
//	rawstr := []string{}
//	for _, r := range raw {
//		rawstr = append(rawstr, fmt.Sprintf("%02X", r))
//	}
//
//	return fmt.Sprintf("%5s %-10s ; %04X %s",
//		instr.Instr.Name, argstr, addr, rawstr,
//	)
//}
