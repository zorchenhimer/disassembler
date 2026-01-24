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

	// Pass 1
	for _, bank := range cfg.Banks {
		//if bank.Name == "bank_01" {
		//	verbose = true
		//}

		//lastoffs := 0
		for offs := bank.Offset; offs < bank.Offset + bank.Size ; {
			if offs >= uint(len(raw)) {
				return fmt.Errorf("offset(%X;%d) past end of input(%X;%d)", offs, offs, len(raw), len(raw))
			}

			decoded := &types.Decoded{
				Type: bank.Type(offs + bank.Address),
			}

			if decoded.Type == types.Range_Code {
				instr := instructions.TryInstr_6502(raw[offs:])
				if instr != nil {
					decoded.Instr = instr
					decoded.Raw = instr.Raw()
					//vb("%04X %#v", offs, decoded)

					switch instr.Instr.AddrMode {
					case types.AddrMode_Accumulator,
						 types.AddrMode_Implied,
						 types.AddrMode_Immediate:
						// nope

					case types.AddrMode_Relative:
						addr := uint(instr.Arg + int(offs + bank.Address))
						//bank.AutoLabels[addr] = types.NewLabel(addr, fmt.Sprintf("L%04X", addr))
						lm.SetLabel(types.NewLabel(addr, fmt.Sprintf("L%04X", addr)))

					default:
						//bank.AutoLabels[uint(instr.Arg)] = types.NewLabel(uint(instr.Arg), fmt.Sprintf("L%04X", instr.Arg))
						lm.SetLabel(types.NewLabel(uint(instr.Arg), fmt.Sprintf("L%04X", instr.Arg)))
					}
				} else {
					decoded.Raw = []byte{raw[offs]}
					decoded.Val = uint(raw[offs])
				}

			} else if decoded.Type == types.Range_Words {
				if uint(len(raw)) <= offs+1 {
					decoded.Raw = []byte{raw[offs]}
					decoded.Val = uint(raw[offs])
				} else {
					decoded.Raw = []byte{raw[offs], raw[offs+1]}
					decoded.Val = uint(raw[offs]) | (uint(raw[offs+1])<<8)
				}
			} else {
				decoded.Raw = []byte{raw[offs]}
				decoded.Val = uint(raw[offs])
			}

			for i := uint(0); i < uint(len(decoded.Raw)); i++ {
				addr := bank.Address+i+(offs-bank.Offset)
				bank.Decoded[addr] = decoded
				//vb("[%04X] %s %#v", addr, decoded.Type, decoded)
			}

			//fmt.Println(decoded.Asm(offs+bank.Address))
			offs += uint(len(decoded.Raw))
		}
	}

	verbose = false
	//return nil

	for _, bank := range cfg.Banks {
		output, err := os.Create(bank.Output)
		if err != nil {
			return err
		}

		for i := bank.Address; i < bank.Address + bank.Size; {
			dec := bank.Decoded[i]
			//if dec == nil {
			//	panic(fmt.Sprintf("[%s] %04X+%04X (%X) how in the fuck?", bank.Name, bank.Address, i, bank.Address + bank.Size))
			//}
			//lbl := bank.Labels[i]
			lbl := lm.GetLabel(i)
			if lbl != nil {
				fmt.Fprintln(output, lbl.Name+":")
			}
			fmt.Fprintln(output, dec.Asm(i, lm))
			i += uint(len(dec.Raw))
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
