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

	for _, win := range windows {
		//wc := &WindowConfig{
		//	Name: win.Name,
		//	Bank: nil,
		//}

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

	return lm
}

//type WindowConfig struct {
//	Name string
//	Bank *Bank
//}

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

	//for _, win := range lm.Windows {
	//	if win.Bank == nil {
	//		continue
	//	}

	//	if win.Bank.Address <= addr && addr < win.Bank.Address + win.Bank.Size {
	//		return win.Bank.Labels[addr]
	//	}
	//}

	//return nil
}

// for autolabels only
func (lm *LabelManager) SetLabel(lbl *types.Label) {
	if lbl == nil {
		return
	}

	_, ok := lm.Global[lbl.Address]
	if ok {
		return
	}

	// Default to global if no window exists for address
	win, ok := lm.WindowAddrs[lbl.Address]
	if !ok || win == "" {
		lm.Global[lbl.Address] = lbl
		return
	}

	bank, ok := lm.Windows[win]
	if !ok || bank == nil {
		return
	}

	//if bank == nil {
	//	panic("[SetLabel] nil bank")
	//}
	// do not overwrite labels
	_, ok = bank.Labels[lbl.Address]
	if ok {
		return
	}
	bank.Labels[lbl.Address] = lbl

	//_, ok = lm.Bank.Labels[lbl.Address]
	//if !ok {
	//	lm.Bank.Labels[lbl.Address] = lbl
	//}

	//for _, win := range lm.Windows {
	//	if win.Bank == nil {
	//		continue
	//	}

	//	if win.Bank.Address <= addr && addr < win.Bank.Address + win.Bank.Size {
	//		_, ok := win.Bank.Labels[addr]
	//		if !ok {
	//			win.Bank.Labels[addr] = lbl
	//			return
	//		}
	//	}
	//}
}

func (lm *LabelManager) GetRange(addr uint) *types.Range {
	//for _, win := range lm.Windows {
	//	if win.Bank == nil {
	//		continue
	//	}

	//	if win.Bank.Address <= addr && addr < win.Bank.Address + win.Bank.Size {
	//		return win.Bank.Ranges[addr]
	//	}
	//}
	//return nil

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

//type WindowManager struct {
//	State map[uint]*WindowState
//}
//
//type WindowState struct {
//	Bank *Bank
//}
//
//func (w *WindowManager) SetLabel(lbl *Label) {
//}
//
//func (w *WindowManager) GetLabel(address uint) *Label {
//}
