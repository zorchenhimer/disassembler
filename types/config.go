package types

import (
	"fmt"
	"strings"
)

type Config struct {
	Global   ConfigGlobal
	Banks    []*Bank
}

func (c *Config) GetInputs() []string {
	files := []string{}
	if c.Global.Input != ""  {
		files = append(files, c.Global.Input)
	}

	for _, bank := range c.Banks {
		if bank.Input != "" && bank.Input != "-" {
			files = append(files, bank.Input)
		}
	}

	return files
}

func (c *Config) Verify() error {
	err := c.Global.verify()
	if err != nil {
		return err
	}

	if len(c.Banks) == 0 {
		return fmt.Errorf("At least one Bank block is required")
	}

	for _, bank := range c.Banks {
		err = bank.verify()
		if err != nil {
			return err
		}

		for _, win := range bank.CfgWindows {
			if win.Address < bank.Address || win.Address >= bank.Address+bank.Size {
				return fmt.Errorf("BankWindow switch at %04X is out of range of bank %s",
					win.Address, bank)
			}

			found := false
			for _, gwin := range c.Global.Windows {
				if win.Window == gwin.Name {
					found = true
				}
			}
			if !found {
				return fmt.Errorf("WindowDef for window %q not found", win.Window)
			}
		}
	}

	for _, win := range c.Global.Windows {
		if win.Init == "" {
			continue
		}

		found := false
		for _, bank := range c.Banks {
			if bank.Name == win.Init {
				found = true
			}
		}

		if !found {
			return fmt.Errorf("Window %s specifies Init bank that does not exist: %s",
				win.Name, win.Init)
		}
	}

	return nil
}

func (c *Config) GoString() string {
	banks := []string{}

	for _, bank := range c.Banks {
		banks = append(banks, fmt.Sprintf("%#v", bank))
	}

	return fmt.Sprintf("{Config Global:%#v Banks:[%s]}",
		c.Global,
		strings.Join(banks, ", "),
	)
}

