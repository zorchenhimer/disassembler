package types

import (
	"fmt"
	"strings"
)

type Range struct {
	Name    string
	Comment string

	Address uint
	//Offset  int
	Size    uint
	End     uint
	Stride  uint // only for byte range

	Type    RangeType
	Display RangeDisplay

	ResolveLabels bool
	RtsLabels     bool
}

func (r *Range) Duplicate() *Range {
	return &Range{
		Name:    r.Name,
		Comment: r.Comment,

		Address: r.Address,
		Size:    r.Size,
		End:     r.End,
		Stride:  r.Stride,

		Type:    r.Type,
		Display: r.Display,

		ResolveLabels: r.ResolveLabels,
		RtsLabels:     r.RtsLabels,
	}
}

type RangeType int

const (
	Range_Code RangeType = iota

	// Block of bytes.
	// Output respects stride.  Defaults to 8.
	Range_Bytes

	// List of two-byte words
	// Output respects stride.  Defaults to 1 if labels, 8 otherwise
	Range_Words
)

func (rt RangeType) String() string {
	switch rt {
	case Range_Bytes:
		return "Range_Bytes"
	case Range_Words:
		return "Range_Words"
	default:
		return "Range_Code"
	}
}

type RangeDisplay int

const (
	Display_Decimal RangeDisplay = iota
	Display_Binary
	Display_Hexadecimal
	Display_Label
)

func (r *Range) Verify(begin, end uint) error {
	if r.Address < begin {
		return fmt.Errorf("Addresses out of range (too low): %#v", r)
	}

	if r.Size == 0 && r.End == 0 {
		return fmt.Errorf("Range missing size or end: %#v", r)
	}

	if r.Size != 0 && r.End != 0 {
		return fmt.Errorf("Range specifies both size and end: %#v", r)
	}

	if r.Size == 0 {
		if r.End < r.Address {
			return fmt.Errorf("Range end before address: %#v", r)
		}
		r.Size = r.End - r.Address+1
	}

	if r.Size + r.Address > end {
		return fmt.Errorf("Addresses out of range (too high): %#v", r)
	}

	if strings.TrimSpace(r.Comment) != "" {
		lines := []string{}
		for _, line := range strings.Split(r.Comment, "\n") {
			lines = append(lines, strings.TrimSpace(line))
		}
		r.Comment = strings.Join(lines, "\n")
	}

	return nil
}
