package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"

	"git.zorchenhimer.com/Zorchenhimer/dasm"
	"git.zorchenhimer.com/Zorchenhimer/dasm/config"
)

type Arguments struct {
	ConfigFile string `arg:"positional,required"`
}

func main() {
	args := &Arguments{}
	arg.MustParse(args)

	err := run(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args *Arguments) error {
	cfg, err := config.ReadFile(args.ConfigFile)
	//fmt.Printf("%#v\n", cfg)
	//fmt.Printf("%#v\n", cfg.Banks)
	if err != nil {
		return err
	}

	err = dasm.FromConfig(cfg)
	if err != nil {
		return err
	}

	return nil
}
