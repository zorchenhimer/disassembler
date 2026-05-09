package instrsbx

import (
	"fmt"
	"strings"

	"git.zorchenhimer.com/Zorchenhimer/dasm/types"
)

type AstNode interface {
	IsNode() bool
	String() string
	DbgString() string
	Address() uint
}

type AstInstruction struct {
	Arguments []AstNode
	Opcode    *Opcode

	lm types.LabelManager
}

func (n *AstInstruction) Address() uint {
	if len(n.Arguments) > 0 {
		return n.Arguments[0].Address()
	}
	return n.Opcode.Address
}

func (n *AstInstruction) String() string {
	//var lastLabel *types.Label
	//if len(n.Arguments) > 0 {
	//	addr := n.Arguments[len(n.Arguments)-1].Address()
	//	lastLabel = n.lm.GetLabel(addr)
	//}

	args := []string{}
	for _, arg := range n.Arguments {
		args = append(args, arg.String())
	}

	str := []string{ n.Opcode.Instr.Name }

	if n.Opcode.Inline != "" {
		switch n.Opcode.Instr.Opcode {
		case 0xB7: // push_var
			return ":"+n.Opcode.Inline
		case 0xB8: // push_word
			return n.Opcode.Inline
		case 0xBB: // push_data
			return fmt.Sprintf("\"%s\"", n.Opcode.Inline)
		case 0xB9: // push_var_indexed
			return fmt.Sprintf("%s[%s]", n.Opcode.Inline, strings.Join(args, " "))

		default:
			if n.Opcode.Instr.Inline == Inline_NullTerm {
				str = append(str, fmt.Sprintf("%q", n.Opcode.Inline))
			} else {
				str = append(str, n.Opcode.Inline)
			}
		}
	}

	if len(args) != 0 {
		str = append(str, strings.Join(args, " "))
	}

	//lbl := ""
	//if lastLabel != nil {
	//	lbl = lastLabel.Name+": "
	//}

	//return lbl+"("+strings.Join(str, " ")+")"
	return "("+strings.Join(str, " ")+")"
}

func (n *AstInstruction) DbgString() string {
	return fmt.Sprintf("%04X: %s", n.Opcode.Address, n.String())
}

func (n *AstInstruction) IsNode() bool { return true }

type AstStackArg struct {
	Data *StackData
	Value uint
}

func (n *AstStackArg) Address() uint {
	return n.Data.Address
}

func (n *AstStackArg) String() string {
	return fmt.Sprintf("%d", n.Value)
}

func (n *AstStackArg) DbgString() string {
	return fmt.Sprintf("%04X: %s", n.Data.Address, n.String())
}
func (n *AstStackArg) IsNode() bool { return true }

// For loops or functions/routines
type AstScope struct {
	Nodes []AstNode
}

func (n *AstScope) Address() uint {
	if len(n.Nodes) == 0 {
		return 0
	}
	return n.Nodes[0].Address()
}

func (n *AstScope) IsNode() bool { return true }

func (n *AstScope) String() string {
	str := []string{}
	for _, node := range n.Nodes {
		str = append(str, node.String())
	}

	return strings.Join(str, "\n")
}

func (n *AstScope) DbgString() string {
	return "[AstScope]"+n.String()
}

// $6103 = $6100 + $61002
/*
	push_var $6100
	push_var $6102
	add
	pop_into $6103

	pop_into $6103 (
		add (
			push_var $6100
			push_var $6102
		)
	)

*/

/*
var stack = []AstNode{
	&AstInstruction{
		Opcode: &Opcode{
			Address: 0x07,
			Instr: Instructions[0xBD], // pop_into
		},
		Arguments: []AstNode{
			&AstInstruction{
				Opcode: Instructions[0xCB], // add
				Arguments: []AstNode{
					&AstInstruction{
						Opcode: Instructions[0xB7], // push_var
						Arguments: []AstNode{
							AstInlineValue{ Value: 0x6100 },
						},
					},
					&AstInstruction{
						Opcode: Instructions[0xB7], // push_var
						Arguments: []AstNode{
							AstInlineValue{ Value: 0x6102 },
						},
					},
				},
			},
		},
	},
}
*/
