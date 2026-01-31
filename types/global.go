package types

import (
	"fmt"
	"strings"
)

type ConfigGlobal struct {
	Input     string
	MlbOutput string

	Comments CommentLevel

	Architecture ArchitectureType

	Windows []*WindowDef
	Labels  []*Label
	Include []string
}

func (g ConfigGlobal) verify() error {
	if strings.TrimSpace(g.Input) == "" {
		return fmt.Errorf("Input missing")
	}

	if g.Architecture != Arch_6502 {
		return fmt.Errorf("Only Architecture 6502 is supported for now.")
	}

	for _, lbl := range g.Labels {
		err := lbl.Verify(0x0000, 0xFFFF)
		if err != nil {
			return fmt.Errorf("Label error: %w", err)
		}
	}

	//for _, win := range g.Windows {
	//	err := win.verify()
	//	if err != nil {
	//		return fmt.Errorf("Windows error: %w", err)
	//	}
	//}

	return nil
}

func (g ConfigGlobal) GoString() string {
	windows := []string{}
	labels := []string{}

	for _, itm := range g.Windows {
		windows = append(windows, fmt.Sprintf("%#v", itm))
	}

	for _, itm := range g.Labels {
		labels = append(labels, fmt.Sprintf("%#v", itm))
	}

	return fmt.Sprintf("{Global Input:%q MlbOutput:%q Comments:%s Architecture:%s Windows:[%s] Labels:[%s] Include:[%s]",
		g.Input,
		g.MlbOutput,
		g.Comments,
		g.Architecture,
		strings.Join(windows, ", "),
		strings.Join(labels, ", "),
		strings.Join(g.Include, ", "),
	)
}

type CommentLevel int

const (
	Comment_None CommentLevel = iota
	Comment_Std
	Comment_Full
)

func (cl CommentLevel) String() string {
	switch cl {
	case Comment_None:
		return "Comment_None"
	case Comment_Std:
		return "Comment_Std"
	case Comment_Full:
		return "Comment_Full"
	}
	return "Comment_UNKNOWN"
}

type ArchitectureType int

const (
	// Offical instructions only
	Arch_6502 ArchitectureType = iota

	// Includes unofficial instructions
	Arch_Full6502

	// StudyBox tape script
	Arch_SbxScript
)

func (arch ArchitectureType) String() string {
	switch arch {
	case Arch_6502:
		return "Arch_6502"
	case Arch_Full6502:
		return "Arch_Full6502"
	case Arch_SbxScript:
		return "Arch_SbxScript"
	}
	return "Arch_UNKNOWN"
}
