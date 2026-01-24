package types

import (
	"fmt"
	"strings"
)

type Label struct {
	Name    string
	Comment string

	Address uint
	//Offset  int
	Size    uint

	ParamSize uint
}

func NewLabel(address uint, name string) *Label {
	return &Label{
		Address:   address,
		Comment:   "",
		Name:      name,
		ParamSize: 0,
		Size:      1,
	}
}

func (lbl *Label) Verify(begin, end uint) error {
	//if lbl.Address < 0 {
	//	return fmt.Errorf("Addresses cannot be negative")
	//}

	if lbl.Address > 0xFFFF {
		return fmt.Errorf("Addresses cannot be larger than $FFFF")
	}

	if lbl.Address < begin {
		return fmt.Errorf("Addresses out of range (too low)")
	}

	if lbl.Size < 0 {
		return fmt.Errorf("Size cannot be negative")
	}

	if lbl.Size + lbl.Address - 1 > 0xFFFF {
		return fmt.Errorf("Size goes beyond $FFFF")
	}

	if lbl.Size + lbl.Address > end {
		return fmt.Errorf("Addresses out of range (too high)")
	}

	if strings.TrimSpace(lbl.Name) == "" && 
	   strings.TrimSpace(lbl.Comment) == "" {
		return fmt.Errorf("Name and Command cannot both be empty")
	}

	if strings.TrimSpace(lbl.Comment) != "" {
		lines := []string{}
		for _, line := range strings.Split(lbl.Comment, "\n") {
			lines = append(lines, strings.TrimSpace(line))
		}
		lbl.Comment = strings.Join(lines, "\n")
	}

	return nil
}

