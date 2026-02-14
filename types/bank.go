package types

import (
	"fmt"
	"strings"
	//"slices"

	//"git.zorchenhimer.com/Zorchenhimer/dasm/instructions"
)

type Bank struct {
	Name   string
	Input  string
	Output string
	NoDasm bool

	Address uint
	Offset  uint
	Size    uint

	CfgLabels  []*Label
	CfgRanges  []*Range
	CfgWindows []*BankWindow

	// key is address
	Labels  map[uint]*Label
	Ranges  map[uint]*Range
	Windows map[uint][]*BankWindow

	Decoded map[uint]AsmLine

	AutoLabels map[uint]*Label
}

func NewBank() *Bank {
	return &Bank{
		Labels:  make(map[uint]*Label),
		Ranges:  make(map[uint]*Range),
		Decoded: make(map[uint]AsmLine),
		Windows: make(map[uint][]*BankWindow),

		CfgLabels:  []*Label{},
		CfgRanges:  []*Range{},
		CfgWindows: []*BankWindow{},

		AutoLabels: make(map[uint]*Label),
	}
}

func (b *Bank) GetName() string {
	return b.Name
}

func (b *Bank) AddrType(address uint) RangeType {
	if rng, ok := b.Ranges[address]; ok {
		return rng.Type
	}
	return Range_Code
}

func (b *Bank) GoString() string {
	labels  := []string{}
	ranges  := []string{}

	for _, itm := range b.CfgLabels {
		labels = append(labels, fmt.Sprintf("%#v", itm))
	}

	for _, itm := range b.CfgRanges {
		ranges = append(ranges, fmt.Sprintf("%#v", itm))
	}

	return fmt.Sprintf("{Bank Name:%q Output:%q Offset:%d Address:%d Size:%d Labels:[%s] Ranges:[%s]}",
		b.Name,
		b.Output,
		b.Offset,
		b.Address,
		b.Size,
		strings.Join(labels, ", "),
		strings.Join(ranges, ", "),
	)
}

func (b *Bank) String() string {
	if b.Name != "" {
		return b.Name
	}

	return fmt.Sprintf("at offset $%X", b.Offset)
}

func (b *Bank) verify() error {
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

	for _, win := range b.CfgWindows {
		b.Windows[win.Address] = append(b.Windows[win.Address], win)
	}

	for _, lbl := range b.CfgLabels {
		for i := uint(0); i < lbl.Size; i++ {
			if _, ok := b.Labels[i+lbl.Address]; ok {
				errs = append(errs, fmt.Errorf("Label overlap at $%04X", i+lbl.Address))
			}
			b.Labels[i+lbl.Address] = lbl
		}
	}

	for _, rng := range b.CfgRanges {
		var err error
		if b.Size == 0 {
			err = rng.Verify(b.Address, b.Address+0x8000)
		} else {
			err = rng.Verify(b.Address, b.Address+b.Size)
		}
		if err != nil {
			errs = append(errs, err)
			continue
		}
		var lbl *Label
		if rng.Size == 0 {
			errs = append(errs, fmt.Errorf("Range with no size: %#v", rng))
		}

		if rng.Name != "" || rng.Comment != "" {
			lbl = &Label{
				Name:         rng.Name,
				CommentBlock: rng.Comment,
				Address:      rng.Address,
				Size:         rng.Size,
			}
		}

		for i := uint(0); i < rng.Size; i++ {
			if thing, ok := b.Ranges[i+rng.Address]; ok {
				errs = append(errs, fmt.Errorf("Range overlap at $%04X: %#v", i+rng.Address, thing))
			}
			b.Ranges[i+rng.Address] = rng
			if lbl != nil {
				if thing, ok := b.Labels[i+rng.Address]; ok {
					errs = append(errs, fmt.Errorf("Label overlap from range at $%04X ($%04X): %#v", rng.Address, rng.Address+i, thing))
				}
				b.Labels[i+rng.Address] = lbl
			}
		}
	}

	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		return fmt.Errorf("Errors encountered")
	}

	return nil
}

func (b *Bank) GetLabel(address uint) *Label {
	if lbl, ok := b.Labels[address]; ok {
		return lbl
	}
	return nil
}

func (b *Bank) SetLabel(lbl *Label) *Label {
	if !(b.Address <= lbl.Address && lbl.Address < b.Address + b.Size) {
		return nil
	}

	if l, ok := b.Labels[lbl.Address]; ok {
		if l.Name == "" {
			l.Name = lbl.Name
		}
		return l
	}

	b.Labels[lbl.Address] = lbl
	return lbl
}

