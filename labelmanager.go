package dasm

import (
	"git.zorchenhimer.com/Zorchenhimer/dasm/types"
)

type LabelManager struct {
	// keys are names
	Windows map[string]*types.Bank
	Banks   map[string]*types.Bank

	// map of addresses to window
	WindowAddrs map[uint]string

	// Global label overrides
	Global map[uint]*types.Label

	defs []*types.WindowDef
}

func NewLabelManager(global []*types.Label, banks []*types.Bank, windows []*types.WindowDef) *LabelManager {
	lm := &LabelManager{
		Windows: make(map[string]*types.Bank),
		Banks:   make(map[string]*types.Bank),
		WindowAddrs: make(map[uint]string),
		Global:  make(map[uint]*types.Label),
	}

	for _, lbl := range global {
		for i := lbl.Address; i < lbl.Address+lbl.Size; i++ {
			lm.Global[i] = lbl
		}
	}

	for _, bank := range banks {
		lm.Banks[bank.Name] = bank
	}

	lm.defs = windows
	lm.Init()

	return lm
}

func (lm *LabelManager) Init() {
	for _, win := range lm.defs {
		lm.Windows[win.Name] = nil

		if win.Init != "" {
			if b, ok := lm.Banks[win.Init]; ok {
				lm.Windows[win.Name] = b
			}
		}

		for i := win.Start; i < win.Start+win.Size; i++ {
			lm.WindowAddrs[i] = win.Name
		}
	}
}

func (lm *LabelManager) SetWindow(winName, bankName string) {
	b, ok := lm.Banks[bankName]
	if !ok {
		return
	}

	_, ok = lm.Windows[winName]
	if !ok {
		return
	}

	lm.Windows[winName] = b
}

func (lm *LabelManager) GetLabel(addr uint) *types.Label {
	lbl, ok := lm.Global[addr]
	if ok {
		return lbl
	}

	win, ok := lm.WindowAddrs[addr]
	if !ok {
		return nil
	}

	bank, ok := lm.Windows[win]
	if ok && bank != nil {
		return bank.Labels[addr]
	}

	return nil
}

// for autolabels only
func (lm *LabelManager) SetLabel(lbl *types.Label) *types.Label {
	if lbl == nil {
		return nil
	}

	gbl, ok := lm.Global[lbl.Address]
	if ok && gbl.Name == "" {
		lm.Global[lbl.Address].Name = lbl.Name
		lm.Global[lbl.Address].References++
		return lm.Global[lbl.Address]
	} else if ok {
		return gbl
	}

	// Default to global if no window exists for address
	win, ok := lm.WindowAddrs[lbl.Address]
	if !ok || win == "" {
		lm.Global[lbl.Address] = lbl
		lm.Global[lbl.Address].References++
		return lm.Global[lbl.Address]
	}

	bank, ok := lm.Windows[win]
	if !ok || bank == nil {
		return nil
	}

	// do not overwrite labels
	// TODO: anon-labels outside of a specific range (window? +-127 bytes?)
	l, ok := bank.Labels[lbl.Address]
	if !ok {
		bank.Labels[lbl.Address] = lbl
		l = lbl
	} else if l.Name == "" {
		l.Name = lbl.Name
	}

	l.References++
	return l
}

func (lm *LabelManager) GetRange(addr uint) *types.Range {
	win, ok := lm.WindowAddrs[addr]
	if !ok {
		return nil
	}

	bank, ok := lm.Banks[win]
	if !ok {
		return nil
	}

	return bank.Ranges[addr]
}
