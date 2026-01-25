package config

import (
	//"fmt"
	//"strings"

	"git.zorchenhimer.com/Zorchenhimer/dasm/types"
)

func defaultRange() *types.Range {
	return &types.Range{
		Size:    0,
		Stride:  8,
		Type:    types.Range_Bytes,
		Display: types.Display_Hexadecimal,
	}
}
