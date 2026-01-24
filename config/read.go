package config

import (
	//"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"git.zorchenhimer.com/Zorchenhimer/dasm/types"
)

func ReadFile(filename string) (*types.Config, error) {
	raw, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Unable to read %s: %w", filename, err)
	}

	return realRead(filepath.Dir(filename), raw, false)
}

func Read(r io.Reader) (*types.Config, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("Unable to read: %w", err)
	}

	return realRead("", raw, false)
}

func realRead(basedir string, raw []byte, nested bool) (*types.Config, error) {
	l, items := newLexer(raw)
	go l.Run()

	cfg, err := newParser(items).Run()
	if err != nil {
		return nil, err
	}

	if nested && len(cfg.Global.Include) > 0 {
		return nil, fmt.Errorf("Nested includes not allowed")
	}

	for _, file := range cfg.Global.Include {
		//inc, err := ReadFile(file)
		fullname := filepath.Join(basedir, file)
		incraw, err := os.ReadFile(fullname)
		if err != nil {
			return nil, fmt.Errorf("Error reading included file %s: %w", file, err)
		}

		inc, err := realRead(filepath.Dir(fullname), incraw, true)
		if err != nil {
			return nil, fmt.Errorf("Error reading included file %s: %w", file, err)
		}

		cfg.Banks = append(cfg.Banks, inc.Banks...)
		cfg.RamBanks = append(cfg.RamBanks, inc.RamBanks...)
	}

	if nested {
		return cfg, nil
	}

	err = cfg.Verify()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
