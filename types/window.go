package types

type BankWindow struct {
	Address uint
	// TODO: input validation with the config
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
	Init  string
}

func (w *WindowDef) verify() error {
	return nil
}
