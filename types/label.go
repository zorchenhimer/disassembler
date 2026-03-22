package types

import (
	"fmt"
	"strings"
)

type LabelManager interface {
	GetLabel(addr uint) *Label
	SetLabel(lbl *Label) *Label
	GetRange(addr uint) *Range

	SetWindow(winName, bankName string)

	// Set the bank that is currently being disassembled.  This should
	// find the window at the address and set the bank in that window.
	ActivateBank(name string, address uint)
}

type Label struct {
	Name    string
	//Comment string

	CommentBlock  string // Block comment starting at the beginning of a line
	CommentInline string // Multi-line comment after the instruction

	Address uint
	//Offset  int
	Size    uint

	ParamSize uint

	References int
}

func NewLabel(address uint, name string) *Label {
	return &Label{
		Address:   address,
		Name:      name,
		ParamSize: 0,
		Size:      1,
	}
}

func (lbl *Label) Verify(begin, end uint) error {
	if lbl.Address > 0xFFFF {
		return fmt.Errorf("Addresses cannot be larger than $FFFF: %#v", lbl)
	}

	if lbl.Address < begin {
		return fmt.Errorf("Addresses out of range (too low): %#v", lbl)
	}

	if lbl.Size < 0 {
		return fmt.Errorf("Size cannot be negative: %#v", lbl)
	}

	if lbl.Size + lbl.Address - 1 > 0xFFFF {
		return fmt.Errorf("Size goes beyond $FFFF: %#v", lbl)
	}

	if lbl.Size + lbl.Address > end {
		return fmt.Errorf("Addresses out of range (too high): %#v", lbl)
	}

	if strings.TrimSpace(lbl.Name) == "" && 
	   strings.TrimSpace(lbl.CommentBlock) == "" &&
	   strings.TrimSpace(lbl.CommentInline) == "" {
		   return fmt.Errorf("Name and Comment cannot both be empty: %#v", lbl)
	}

	return nil
}

