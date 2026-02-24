package types

type BankWindow struct {
	Address uint
	Window  string
	Bank    string
}

func (bw BankWindow) String() string {
	return bw.Window
}

type WindowDef struct {
	Name  string
	Start uint
	Size  uint
	Init  string
}

func (w *WindowDef) verify() error {
	return nil
}
