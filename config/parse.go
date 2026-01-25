package config

import (
	"fmt"
	"strconv"
	"strings"

	"git.zorchenhimer.com/Zorchenhimer/dasm/types"
)

type parser struct {
	config *types.Config
	items chan lexItem
	verbose bool
}

func newParser(items chan lexItem) *parser {
	return &parser{
		config: &types.Config{},
		items: items,
		//verbose: true,
	}
}

func (p *parser) next() lexItem {
	for {
		select {
		case itm := <- p.items:
			if p.verbose {
				fmt.Printf("%#v\n", itm)
			}
			if itm.typ == lex_Space || itm.typ == lex_Comment {
				continue
			}
			return itm
		}
	}
	return EOF_Item
}

func (p *parser) Run() (*types.Config, error) {
	var err error

	for {
		itm := p.next()
		if itm.typ == lex_EOF {
			break
		}

		switch itm.typ {
		case lex_Comment, lex_Space:
			continue

		case lex_Ident:
			// yup

		default:
			return p.config, parseError(itm, "Invalid config: %#v", itm)
		}

		switch strings.ToLower(itm.val) {
		case "global":
			err = p.parseGlobal()
		case "bank":
			err = p.parseBank()
		case "rambank":
			err = p.parseRamBank()
		default:
			err = parseError(itm, "Invalid block: %s", itm.val)
		}

		if err != nil {
			return p.config, err
		}
	}

	return p.config, nil
}

func (p *parser) parseGlobal() error {
	itm := p.next()
	if itm.typ != lex_OpenBracket {
		return parseError(itm, "Missing open bracket after Global")
	}

	for {
		itm = p.next()
		if itm.typ == lex_Comment || itm.typ == lex_Semicolon {
			continue
		}

		if itm.typ == lex_CloseBracket {
			break
		}

		if itm.typ == lex_EOF {
			return fmt.Errorf("Early EOF")
		}

		if itm.typ != lex_Ident {
			return parseError(itm, "Expected IDENT")
		}

		prev := itm
		switch strings.ToLower(itm.val) {
		case "input", "mlboutput", "include":
			itm = p.next()
			if itm.typ != lex_String {
				return parseError(itm, "%s expects a string", prev.val)
			}

			switch strings.ToLower(prev.val) {
			case "input":
				p.config.Global.Input = itm.val
			case "mlboutput":
				p.config.Global.MlbOutput = itm.val
			case "include":
				p.config.Global.Include = append(p.config.Global.Include, itm.val)
			}

		case "architecture":
			itm = p.next()
			if itm.typ != lex_Ident && itm.typ != lex_Number {
				return parseError(itm, "%s expects an IDENT; got %s", prev.val, itm.typ)
			}

			switch strings.ToLower(itm.val) {
			case "6502":
				p.config.Global.Architecture = types.Arch_6502
			case "full6502":
				p.config.Global.Architecture = types.Arch_Full6502
			case "sbx":
				p.config.Global.Architecture = types.Arch_SbxScript
			default:
				return parseError(itm, "Invalid architecture: %s", itm.val)
			}

		case "comments":
			itm = p.next()
			if itm.typ != lex_Ident {
				return parseError(itm, "%s expects an IDENT", prev.val)
			}

			switch strings.ToLower(itm.val) {
			case "none":
				p.config.Global.Comments = types.Comment_None
			case "standard", "std":
				p.config.Global.Comments = types.Comment_Std
			case "full":
				p.config.Global.Comments = types.Comment_Full
			default:
				return parseError(itm, "Invalid comment: %s", itm.val)
			}

		case "windows":
			vals, err := p.parseWindowDefs()
			if err != nil {
				return err
			}
			p.config.Global.Windows = vals

		case "labels":
			vals, err := p.parseLabels()
			if err != nil {
				return err
			}
			p.config.Global.Labels = vals

		default:
			return parseError(itm, "Invalid IDENT: %s", itm.val)
		}
	}

	return nil
}

func (p *parser) parseBank() error {
	itm := p.next()
	if itm.typ != lex_OpenBracket {
		return parseError(itm, "Missing open bracket after Bank")
	}

	//p.verbose = true

	bank := types.NewBank()
	for {
		itm = p.next()
		if itm.typ == lex_Comment || itm.typ == lex_Semicolon {
			continue
		}

		if itm.typ == lex_CloseBracket {
			break
		}

		if itm.typ == lex_EOF {
			return fmt.Errorf("Early EOF")
		}

		if itm.typ != lex_Ident {
			return parseError(itm, "Expected IDENT")
		}

		prev := itm
		switch strings.ToLower(prev.val) {
		case "name", "output":
			itm = p.next()
			if itm.typ != lex_String {
				return parseError(itm, "%s expects a string", prev.val)
			}

			switch strings.ToLower(prev.val) {
			case "name":
				bank.Name = itm.val
			case "output":
				bank.Output = itm.val
			}

		case "address", "offset", "size":
			itm = p.next()
			if itm.typ != lex_Number {
				return parseError(itm, "%s requires lex_Number", prev.val)
			}

			num, err := parseNumber(itm.val)
			if err != nil {
				return err
			}

			switch strings.ToLower(prev.val) {
			case "address":
				bank.Address = num
			case "offset":
				bank.Offset = num
			case "size":
				bank.Size = num
			}

		case "labels":
			vals, err := p.parseLabels()
			if err != nil {
				return err
			}
			bank.CfgLabels = vals

		case "windows":
			vals, err := p.parseBankWindow()
			if err != nil {
				return err
			}
			bank.CfgWindows = vals

		case "ranges":
			vals, err := p.parseRanges()
			if err != nil {
				return err
			}
			bank.CfgRanges = vals

		default:
			return parseError(itm, "Invalid item: %s", itm.val)
		}
	}

	p.config.Banks = append(p.config.Banks, bank)
	return nil
}

func (p *parser) parseBankWindow() ([]*types.BankWindow, error) {
	itm := p.next()
	if itm.typ != lex_OpenSquare {
		return nil, parseError(itm, "Windows expects a list")
	}

	list := []*types.BankWindow{}
	for {
		itm = p.next()
		if itm.typ == lex_CloseSquare {
			break
		}

		if itm.typ != lex_OpenBracket {
			return nil, parseError(itm, "Expected lex_OpenBracket, got %#v", itm)
		}

		win := &types.BankWindow{}
		for {
			itm = p.next()
			if itm.typ == lex_Semicolon {
				continue
			}
			if itm.typ == lex_CloseBracket {
				break
			}

			if itm.typ != lex_Ident {
				return nil, parseError(itm, "Expected lex_Ident, got %#v", itm)
			}

			key := itm
			val := p.next()
			switch strings.ToLower(itm.val) {
			case "bank", "window":
				if val.typ != lex_String {
					return nil, parseError(itm, "%s requires lex_String", key.val)
				}

				if strings.ToLower(key.val) == "bank" {
					win.Bank = val.val
				} else {
					win.Window = val.val
				}

			case "address":
				if val.typ != lex_Number {
					return nil, parseError(itm, "%s requires lex_Number", key.val)
				}

				num, err := parseNumber(val.val)
				if err != nil {
					return nil, err
				}
				win.Address = num

			default:
				return nil, parseError(itm, "Invalid item: %s", itm.val)
			}
		}
		list = append(list, win)
	}

	return list, nil
}

func (p *parser) parseRamBank() error {
	itm := p.next()
	if itm.typ != lex_OpenBracket {
		return parseError(itm, "Missing open bracket after RamBank")
	}

	bank := &types.RamBank{}
	for {
		itm = p.next()
		if itm.typ == lex_Semicolon {
			continue
		}

		if itm.typ == lex_CloseBracket {
			break
		}

		if itm.typ == lex_EOF {
			return fmt.Errorf("Early EOF")
		}

		if itm.typ != lex_Ident {
			return parseError(itm, "Expected IDENT")
		}

		prev := itm
		switch strings.ToLower(prev.val) {
		case "name", "output":
			itm = p.next()
			if itm.typ != lex_String {
				return parseError(itm, "%s expects a string", prev.val)
			}

			switch strings.ToLower(prev.val) {
			case "name":
				bank.Name = itm.val
			case "output":
				bank.Output = itm.val
			}

		case "address", "offset", "size":
			itm = p.next()
			if itm.typ != lex_Number {
				return parseError(itm, "%s requires lex_Number", prev.val)
			}

			num, err := parseNumber(itm.val)
			if err != nil {
				return err
			}

			switch strings.ToLower(prev.val) {
			case "address":
				bank.Address = num
			case "offset":
				bank.Offset = num
			case "size":
				bank.Size = num
			}

		case "labels":
			vals, err := p.parseLabels()
			if err != nil {
				return err
			}
			bank.CfgLabels = vals
		}
	}

	p.config.RamBanks = append(p.config.RamBanks, bank)
	return nil
}

func (p *parser) parseWindowDefs() ([]*types.WindowDef, error) {
	itm := p.next()
	if itm.typ != lex_OpenSquare {
		return nil, parseError(itm, "Windows expects a list")
	}

	list := []*types.WindowDef{}

	for {
		itm = p.next()
		if itm.typ == lex_CloseSquare {
			break
		}

		if itm.typ != lex_OpenBracket {
			return nil, parseError(itm, "Expected lex_OpenBracket, got %#v", itm)
		}

		win := &types.WindowDef{}
		for {
			itm = p.next()
			if itm.typ == lex_Semicolon {
				continue
			}
			if itm.typ == lex_CloseBracket {
				break
			}

			if itm.typ != lex_Ident {
				return nil, parseError(itm, "Expected lex_Ident, got %#v", itm)
			}

			key := itm
			val := p.next()
			switch strings.ToLower(itm.val) {
			case "name", "init":
				if val.typ != lex_String {
					return nil, parseError(itm, "%s requires lex_String", key.val)
				}

				if strings.ToLower(key.val) == "name" {
					win.Name = val.val
				} else {
					win.Init = val.val
				}

			case "start", "size":
				if val.typ != lex_Number {
					return nil, parseError(itm, "%s requires lex_Number", key.val)
				}

				num, err := parseNumber(val.val)
				if err != nil {
					return nil, err
				}

				if strings.ToLower(key.val) == "start" {
					win.Start = num
				} else {
					win.Size = num
				}

			case "type":
				if val.typ != lex_Ident {
					return nil, parseError(itm, "%s requires lex_Ident", key.val)
				}

				switch strings.ToLower(val.val) {
				case "ram":
					win.Type = types.Window_Ram
				case "rom":
					win.Type = types.Window_Rom
				default:
					return nil, parseError(itm, "Invalid window type: %s", val.val)
				}

			default:
				return nil, parseError(itm, "Invalid item: %#v", key)
			}
		}

		list = append(list, win)
	}

	return list, nil
}

func (p *parser) parseLabels() ([]*types.Label, error) {
	itm := p.next()
	if itm.typ != lex_OpenSquare {
		return nil, parseError(itm, "Labels expects a list")
	}

	list := []*types.Label{}

	for {
		itm = p.next()
		if itm.typ == lex_CloseSquare {
			break
		}

		if itm.typ != lex_OpenBracket {
			return nil, parseError(itm, "Expected lex_OpenBracket, got %#v", itm)
		}

		lbl := &types.Label{Size:1}
		for {
			itm = p.next()
			if itm.typ == lex_Semicolon {
				continue
			}
			if itm.typ == lex_CloseBracket {
				break
			}

			if itm.typ != lex_Ident {
				return nil, parseError(itm, "Expected lex_Ident, got %#v", itm)
			}

			key := itm
			val := p.next()
			switch strings.ToLower(itm.val) {
			case "address", "size", "paramsize":
				if val.typ != lex_Number {
					return nil, parseError(itm, "%s requires lex_Number", key.val)
				}

				num, err := parseNumber(val.val)
				if err != nil {
					return nil, err
				}

				switch strings.ToLower(key.val) {
				case "address":
					lbl.Address = num
				case "size":
					lbl.Size = num
				case "paramsize":
					lbl.ParamSize = num
				}

			case "name", "comment":
				if val.typ != lex_String {
					return nil, parseError(itm, "%s requires lex_String", key.val)
				}

				if strings.ToLower(key.val) == "name" {
					lbl.Name = val.val
				} else {
					lbl.Comment = val.val
				}

			default:
				return nil, parseError(itm, "Invalid item: %#v", key)
			}
		}

		list = append(list, lbl)
	}
	return list, nil
}

func (p *parser) parseRanges() ([]*types.Range, error) {
	itm := p.next()
	if itm.typ != lex_OpenSquare {
		return nil, parseError(itm, "Labels expects a list")
	}

	list := []*types.Range{}

	for {
		itm = p.next()
		if itm.typ == lex_CloseSquare {
			break
		}

		if itm.typ != lex_OpenBracket {
			return nil, parseError(itm, "Expected lex_OpenBracket, got %#v", itm)
		}

		rng := defaultRange()
		for {
			itm = p.next()
			if itm.typ == lex_Semicolon {
				continue
			}
			if itm.typ == lex_CloseBracket {
				break
			}

			if itm.typ != lex_Ident {
				return nil, parseError(itm, "Expected lex_Ident, got %#v", itm)
			}

			key := itm
			val := p.next()
			switch strings.ToLower(itm.val) {
			case "address", "size", "stride", "end":
				if val.typ != lex_Number {
					return nil, parseError(itm, "%s requires lex_Number", key.val)
				}

				num, err := parseNumber(val.val)
				if err != nil {
					return nil, err
				}

				switch strings.ToLower(key.val) {
				case "address":
					rng.Address = num
				case "size":
					rng.Size = num
				case "stride":
					rng.Stride = num
				case "end":
					rng.End = num
				}

			case "type":
				if val.typ != lex_Ident {
					return nil, parseError(itm, "%s requires lex_Ident", key.val)
				}

				switch strings.ToLower(val.val) {
				case "bytes":
					rng.Type = types.Range_Bytes
				case "code":
					rng.Type = types.Range_Code
				case "words":
					rng.Type = types.Range_Words
				default:
					return nil, parseError(itm, "Invalid range type: %s", val.val)
				}

			case "display":
				if val.typ != lex_Ident {
					return nil, parseError(itm, "%s requires lex_Ident", key.val)
				}

				switch strings.ToLower(val.val) {
				case "bin", "binary":
					rng.Display = types.Display_Binary
				case "dec", "decimal":
					rng.Display = types.Display_Decimal
				case "hex", "hexadecimal":
					rng.Display = types.Display_Hexadecimal
				case "label":
					rng.Display = types.Display_Label
				default:
					return nil, parseError(itm, "Invalid range display: %s", val.val)
				}

			case "resolvelabels", "rtslabels":
				if val.typ != lex_Ident {
					return nil, parseError(itm, "%s requires lex_Ident", key.val)
				}

				b := false
				switch strings.ToLower(val.val) {
				case "true":
					b = true
				case "false":
					b = false
				default:
					return nil, parseError(itm, "%s reqires true/false, got %s", key.val, val.val)
				}

				if strings.ToLower(key.val) == "resolvelabels" {
					rng.ResolveLabels = b
				} else {
					rng.RtsLabels = b
				}

			case "name", "comment":
				if val.typ != lex_String {
					return nil, parseError(itm, "%s requires lex_String", key.val)
				}

				if strings.ToLower(key.val) == "name" {
					rng.Name = val.val
				} else {
					rng.Comment = val.val
				}
			}
		}

		list = append(list, rng)
	}

	//fmt.Printf("parsed %d ranges\n", len(list))
	return list, nil
}

func parseNumber(raw string) (uint, error) {
	if strings.HasPrefix(raw, "$") {
		raw = "0x"+raw[1:]
	} else if strings.HasPrefix(raw, "%") {
		raw = "0b"+raw[1:]
	}

	val, err := strconv.ParseUint(raw, 0, 32)
	if err != nil {
		return 0, err
	}

	return uint(val), nil
}

func parseError(itm lexItem, format string, args ...any) error {
	return fmt.Errorf("[%d:%d] %s", itm.line, itm.col, fmt.Sprintf(format, args...))
}
