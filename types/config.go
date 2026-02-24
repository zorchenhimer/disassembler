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

	hasBankWindows := false
	hasGlobalWindows := len(c.Global.Windows) > 0

	for _, bank := range c.Banks {
		err = bank.verify()
		if err != nil {
			return err
		}

		if bank.Input == "" && c.Global.Input == "" {
			return fmt.Errorf("Bank missing input file")
		} else if bank.Input == "" && c.Global.Input != "" {
			bank.Input = c.Global.Input
		}

		hasBankWindows = len(bank.CfgWindows) > 0
		if hasGlobalWindows && len(bank.CfgWindows) == 0 && !bank.NoDasm {
			// don't warn if a global window init's to this bank.
			init := false
			for _, gwin := range c.Global.Windows {
				if gwin.Init == bank.Name {
					init = true
				}
			}
			if !init {
				fmt.Printf("Warning: Bank %s has no window definitions\n", bank.Name)
			}
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

	if hasBankWindows && !hasGlobalWindows {
		return fmt.Errorf("Bank windows defined without any Global definitions")
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

	// No windows defined anywhere.  Setup some so defined labels work.
	if !hasBankWindows && !hasGlobalWindows {
		fmt.Println("Warning: No windows defined. Creating auto window.")

		// Find bank address bounds.  We want to create a window that can hold all
		// of the banks.  A window that has more space than a given bank shouldn't
		// pose any issue.
		low  := uint(0xFFFF)
		high := uint(0x0000)
		for _, bank := range c.Banks {
			if bank.Address < low {
				low = bank.Address
			}
			if bank.Address + bank.Size - 1 > high {
				high = bank.Address + bank.Size - 1
			}
		}

		c.Global.Windows = append(c.Global.Windows, &WindowDef{
			Name: "auto_window",
			Start: low,
			Size: high - low + 1,
		})

		//fmt.Printf("%#v\n", c.Global.Windows[0])

		// Create a single window entry in each bank at the start address to ensure
		// it's loaded when disassembling.
		for _, bank := range c.Banks {
			bank.Windows[bank.Address] = append(bank.Windows[bank.Address], &BankWindow{
				Address: bank.Address,
				Window:  "auto_window",
				Bank:    bank.Name,
			})
		}

		//for _, bank := range c.Banks {
		//	fmt.Printf("%#v\n", bank.Windows[bank.Address][0])
		//}
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

