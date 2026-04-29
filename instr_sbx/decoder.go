package instrsbx

import (
	"fmt"
	"os"
	"slices"

	"git.zorchenhimer.com/Zorchenhimer/dasm/types"
)

type decoder struct {
	lm types.LabelManager
	autoVars bool
	bankAddr uint
	bankSize uint

	ast   []AstNode // output
	stack *Stack[AstNode] // working vals
}

func NewDecoder(lm types.LabelManager, autoVars bool) types.Decoder {
	return &decoder{
		lm: lm,
		autoVars: autoVars,

		ast:   []AstNode{},
		stack: NewStack[AstNode](),
	}
}

func (dec *decoder) InstrNames() []string {
	names := []string{}
	for _, instr := range Instructions {
		names = append(names, instr.StatName())
	}
	return names
}

// returned line is the opcode an any inline values.  stack arguments are elsewhere.
func (dec *decoder) TryInstr(addr uint, raw []byte) types.AsmLine {
	if len(raw) == 0 {
		return nil
	}

	if addr == 0x6000 {
		return &StackData{
			Address: addr,
			Values: []uint{uint(raw[0]) | (uint(raw[1])<<8)},
			Type: types.Range_Words,
		}
	}

	// value high enough for opcode?
	if raw[0] & 0x80 != 0x80 {
		ret := &StackData{
			Address: addr,
			Values:  []uint{uint(raw[0])},
		}

		dec.stack.Push(&AstStackArg{
			Data: ret,
			Value: ret.Values[0],
		})

		return ret
	}

	instr, ok := Instructions[raw[0]]
	if !ok {
		panic(fmt.Sprintf("Missing SBX Script opcode %02X\n", raw[0]))
	}

	op := &Opcode{
		Address: addr,
		Instr: instr,
	}

	switch op.Instr.Inline {
	case Inline_NullTerm:
		found_end := false
		for i := 1; i < len(raw[1:]) && i < 33; i++ {
			op.RawInline = append(op.RawInline, raw[i])
			if raw[i] == 0x00 {
				found_end = true
				break
			}
		}
		if !found_end {
			fmt.Printf("no end found for %s at %04X\n",
				op.Instr.Name, op.Address)
			return nil
		}

	case Inline_Word, Inline_Label:
		if len(raw) < 3 {
			fmt.Printf("Not enough bytes for inline word for %s at %04X\n",
				op.Instr.Name, op.Address)
			return nil
		}

		op.RawInline = raw[1:3]
		if op.Instr.Inline == Inline_Label {
			lbladdr := uint(raw[1]) | (uint(raw[2])<<8)
			dec.lm.SetLabel(types.NewLabel(lbladdr, fmt.Sprintf("L%04X", lbladdr)))
		}
		//op.vars = append(op.vars, lbladdr)

	case Inline_Count:
		// needs at least one count byte and and one inline word
		if len(raw) < 4 {
			fmt.Printf("Not enough bytes for inline word list for %s at %04X\n",
				op.Instr.Name, op.Address)
			return nil
		}

		count := int(raw[1])
		if len(raw[2:]) < count*2 {
			fmt.Printf("Not enough bytes for %d inline words for %s at %04X\n",
				count, op.Instr.Name, op.Address)
			return nil
		}
		op.RawInline = raw[1:(count*2)+2]
	}

	node := &AstInstruction{ Opcode: op, lm: dec.lm }
	for i := 0; i < op.Instr.WordArgs; i++ {
		n, err := dec.stack.Pop()
		if err != nil {
			fmt.Printf("Opcode pops too much off stack: %#v\n", op)
			return op
		}

		node.Arguments = append(node.Arguments, n)
	}

	for i := 0; i < op.Instr.StringArgs; i++ {
		n, err := dec.stack.Pop()
		if err != nil {
			fmt.Printf("Opcode pops too much off stack: %#v\n", op)
			return op
		}

		node.Arguments = append(node.Arguments, n)
	}

	if len(node.Arguments) > 0 {
		slices.Reverse(node.Arguments)
	}

	if op.Instr.ReturnWord || op.Instr.ReturnString {
		dec.stack.Push(node)
	} else {
		dec.ast = append(dec.ast, node)
	}

	op.Prep(dec.lm)
	return op
}

func (dec *decoder) NewData(addr uint, raw []byte, stride int, display types.RangeDisplay, rngType types.RangeType, rtsLabels bool) types.AsmLine {
	sd := &StackData{
		Address: addr,
		Values: []uint{},
		Display: display,
		Stride: stride,
		Type: rngType,
	}

	if sd.Type == types.Range_Words || sd.Type == types.Range_Addresses {
		for i := 0; i < len(raw); i+=2 {
			if i+1 >= len(raw) {
				break
			}
			sd.Values = append(sd.Values, uint(raw[i]) | (uint(raw[i+1]) << 8))
		}
	} else {
		for i := 0; i < len(raw); i++ {
			sd.Values = append(sd.Values, uint(raw[i]))
		}
	}

	return sd
}

func (dec *decoder) SetBank(addr uint, size uint) {
	dec.bankAddr = addr
	dec.bankSize = size
}

func (dec *decoder) DumpAst(filename string) error {
	output, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer output.Close()

	for _, node := range dec.ast {
		lbl := dec.lm.GetLabel(node.Address())
		lblStr := ""
		if lbl != nil {
			lblStr = "\n"+lbl.Name+": "
		}
		fmt.Fprintf(output, "%s%s\n", lblStr, node.String())
	}

	st := dec.stack.Array()

	if len(st) > 0 {
		fmt.Println("left over stack in", filename)
		fmt.Fprintln(output, "\nleft over stack:")
		for _, node := range st {
			fmt.Fprintln(output, " ", node.DbgString())
		}
	}

	//fmt.Printf("stack: %#v\n", dec.stack)
	//fmt.Println("stack.bottom:", dec.stack.bottom)
	dec.ast = []AstNode{}
	dec.stack = NewStack[AstNode]()
	return nil
}
