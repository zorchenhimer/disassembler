package stats

import (
	"os"
	"fmt"
	"sort"
	"io"
)

func (s *Set) WriteToFile(filename string) error {
	output, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer output.Close()

	fmt.Fprintln(output, "Global")
	err = writeGroup(output, s.Global, 4)
	if err != nil {
		return err
	}
	fmt.Fprintln(output, "")

	for name, group := range s.Groups {
		fmt.Fprintln(output, name)
		err = writeGroup(output, group, 4)
		if err != nil {
			return err
		}
		fmt.Fprintln(output, "")
	}
	return nil

	//keys := []string{}
	//width := 0
	//for key, _ := range s.Global {
	//	keys = append(keys, key)
	//	if len(key) > width {
	//		width = len(key)
	//	}
	//}
	//sort.Strings(keys)
	//width += 3

	//for _, key := range keys {
	//	fmt.Fprintf(output, "%*s %d\n",
	//		width*-1, key, s.Global[key])
	//}

	//for name, group := range s.Groups {
	//	fmt.Fprintln(name)
	//	keys = []string{}
	//	width = 0
	//	for key, _ := range group {
	//		keys = append(keys, key)
	//	}
	//	sort.Strings(keys)
	//	width += 3
	//}

	//return nil
}

func writeGroup(writer io.Writer, group Group, indent int) error {
	keys := []string{}
	width := 0
	for key, _ := range group {
		keys = append(keys, key)
		if len(key) > width {
			width = len(key)
		}
	}
	sort.Strings(keys)
	width += 3

	for _, key := range keys {
		_, err := fmt.Fprintf(writer, "%*s%*s %d\n",
			indent, "", width*-1, key, group[key])
		if err != nil {
			return err
		}
	}

	return nil
}
