package types

import (
	"fmt"
	"strings"
	//"slices"

	//"git.zorchenhimer.com/Zorchenhimer/dasm/instructions"
)

type Bank struct {
	Name   string
	Output string

	Address uint
	Offset  uint
	Size    uint

	CfgLabels  []*Label
	CfgRanges  []*Range
	CfgWindows []*BankWindow

	//Address uint
	//Offset  uint

	// key is address
	Labels map[uint]*Label
	Ranges map[uint]*Range

	Decoded map[uint]*Decoded

	AutoLabels map[uint]*Label

	// currently mapped windows.  key is address
	// (expanded like lables and ranges)
	Windows map[uint]*Bank
}

func NewBank() *Bank {
	return &Bank{
		Labels:  make(map[uint]*Label),
		Ranges:  make(map[uint]*Range),
		Decoded: make(map[uint]*Decoded),

		CfgLabels:  []*Label{},
		CfgRanges:  []*Range{},
		CfgWindows: []*BankWindow{},

		AutoLabels: make(map[uint]*Label),
	}
}

//func (b *Bank) Asm(extraLabels []*Label) string {
//	addrs := []uint{}
//	for addr, _ := range b.Decoded {
//		addrs = append(addrs, addr)
//	}
//	slices.Sort(addrs)
//
//	builder := &strings.Builder{}
//	for _, addr := range addrs {
//		fmt.Fprintln(builder, b.Decoded[i].Asm(addr))
//	}
//
//	return builder.String()
//}

func (b *Bank) Label(address uint) *Label {
	if lbl, ok := b.Labels[address]; ok {
		return lbl
	}
	return nil
}

func (b *Bank) Type(address uint) RangeType {
	if rng, ok := b.Ranges[address]; ok {
		return rng.Type
	}
	return Range_Code
}

func (b *Bank) GoString() string {
	labels  := []string{}
	ranges  := []string{}
	windows := []string{}

	for _, itm := range b.CfgLabels {
		labels = append(labels, fmt.Sprintf("%#v", itm))
	}

	for _, itm := range b.CfgRanges {
		ranges = append(ranges, fmt.Sprintf("%#v", itm))
	}

	for _, itm := range b.CfgWindows {
		windows = append(windows, fmt.Sprintf("%#v", itm))
	}

	return fmt.Sprintf("{Bank Name:%q Output:%q Offset:%d Address:%d Size:%d Labels:[%s] Ranges:[%s] Windows:[%s]}",
		b.Name,
		b.Output,
		b.Offset,
		b.Address,
		b.Size,
		strings.Join(labels, ", "),
		strings.Join(ranges, ", "),
		strings.Join(windows, ", "),
	)
}

func (b *Bank) String() string {
	if b.Name != "" {
		return b.Name
	}

	return fmt.Sprintf("at offset $%X", b.Offset)
}

// TODO: figure out how the windows will be looked up on the consumption end.
func (b *Bank) setupWindows(defs []*WindowDef, banks []*Bank) error {
	//defNames := make(map[string]*WindowDef)
	//for _, def := range defs {
	//	defNames[def.Name] = def
	//}

	//for _, swap := range b.CfgWindows {
	//	for _, bank := range banks {
	//		if swap.Bank == bank.Name {
	//			swap.bnk = bank
	//		}
	//	}
	//}

	//for _, swap := range b.CfgWindows {
	//	b.Windows[swap.Address] = swap.bnk
	//}
	return nil
}

func (b *Bank) verify() error {
	if strings.TrimSpace(b.Output) == "" {
		return fmt.Errorf("Output missing")
	}

	if b.Address < 0 {
		return fmt.Errorf("Address cannot be negative")
	}

	if b.Offset < 0 {
		return fmt.Errorf("Offset cannot be negative")
	}

	if b.Size < 0 {
		return fmt.Errorf("Size cannot be negative")
	}

	if b.Address + b.Size - 1 > 0xFFFF {
		return fmt.Errorf("Size goes beyond $FFFF: $%04X + $%04X = $%04X",b.Address, b.Size, b.Address + b.Size)
	}

	for _, lbl := range b.Labels {
		err := lbl.Verify(b.Address, b.Address+b.Size)
		if err != nil {
			return fmt.Errorf("Label error: %w", err)
		}
	}

	for _, rng := range b.Ranges {
		err := rng.Verify(b.Address, b.Address+b.Size)
		if err != nil {
			return fmt.Errorf("Range error: %w", err)
		}
	}

	errs := []error{}

	for _, lbl := range b.CfgLabels {
		for i := uint(0); i < lbl.Size; i++ {
			if _, ok := b.Labels[i+lbl.Address]; ok {
				errs = append(errs, fmt.Errorf("Label overlap at $%04X", i+lbl.Address))
			}
			b.Labels[i+lbl.Address] = lbl
		}
	}

	for _, rng := range b.CfgRanges {
		for i := uint(0); i < rng.Size; i++ {
			if rng.Name != "" || rng.Comment != "" {
				if _, ok := b.Labels[rng.Address]; ok {
					errs = append(errs, fmt.Errorf("Label overlap from range at $%04X", i+rng.Address))
				}

				b.Labels[rng.Address] = &Label{
					Name:    rng.Name,
					Comment: rng.Comment,
					Address: rng.Address,
					Size:    rng.Size,
				}
			}

			if _, ok := b.Ranges[i+rng.Address]; ok {
				errs = append(errs, fmt.Errorf("Range overlap at $%04X", i+rng.Address))
			}
			b.Ranges[i+rng.Address] = rng
		}
	}

	return nil
}

func (b *Bank) GetLabel(address uint) *Label {
	if b.Address <= address && address < b.Address + b.Size {
		if lbl, ok := b.Labels[address]; ok {
			return lbl
		}
		return nil
	}
	return nil
}

type RamBank struct {
	Name   string
	Output string

	Address uint
	Size    uint
	Offset  uint

	CfgLabels []*Label
}

func (b *RamBank) verify() error {
	return nil
}
