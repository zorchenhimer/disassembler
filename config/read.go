package config

import (
	//"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"git.zorchenhimer.com/Zorchenhimer/dasm/types"
)

func ReadFile(filename string) (*types.Config, error) {
	raw, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Unable to read %s: %w", filename, err)
	}

	cfg, err := realRead(filepath.Dir(filename), raw, false)
	if err != nil {
		return cfg, fmt.Errorf("%s: %w", filename, err)
	}
	return cfg, nil
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

	cfg, err := newParser(items).Run(nested)
	if err != nil {
		return nil, err
	}

	if nested && len(cfg.Global.Include) > 0 {
		return nil, fmt.Errorf("Nested includes not allowed")
	}

	included := []string{}
	for _, item := range cfg.Global.Include {
		item = filepath.Join(basedir, item)
		if strings.HasSuffix(item, string(os.PathSeparator)) {
			item += "*.cfg"
		}

		files, err := filepath.Glob(item)
		if err != nil {
			return nil, err
		}
		included = append(included, files...)
	}

	for _, file := range included {
		fmt.Println(file)
		incraw, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", file, err)
		}

		inc, err := realRead(filepath.Dir(file), incraw, true)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", file, err)
		}

		cfg.Banks = append(cfg.Banks, inc.Banks...)
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
