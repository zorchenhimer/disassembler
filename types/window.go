package types

type BankWindow struct {
	Address uint
	Window  string
	Bank    string

	bnk *Bank
}

func (bw BankWindow) String() string {
	return bw.Window
}

type WindowDef struct {
	Name  string
	Start uint
	Size  uint
	Type  WindowType
	Init  string
}

func (w *WindowDef) verify() error {
	return nil
}

type WindowType int

const (
	Window_Rom WindowType = iota
	Window_Ram
)
