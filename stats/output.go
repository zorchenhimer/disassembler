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
	err = writeGroup(output, s.Global, 4, false)
	if err != nil {
		return err
	}
	fmt.Fprintln(output, "")

	for name, group := range s.Groups {
		fmt.Fprintln(output, name)
		err = writeGroup(output, group, 4, true)
		if err != nil {
			return err
		}
		fmt.Fprintln(output, "")
	}
	return nil
}

func writeGroup(writer io.Writer, group *Group, indent int, onlyUsed bool) error {
	keys := []string{}
	width := 0
	for key, _ := range group.Items {
		keys = append(keys, key)
		if len(key) > width {
			width = len(key)
		}
	}
	sort.Strings(keys)
	width += 3

	fmt.Fprintf(writer, "%*sTotal Size: %d\n", indent, "", group.Size)

	unused := []string{}
	for _, key := range keys {
		str := fmt.Sprintf("%*s%*s %d",
			indent, "", width*-1, key, group.Items[key])
		if group.Items[key] <= 0 {
			unused = append(unused, str)
		} else {
			_, err := fmt.Fprintln(writer, str)
			if err != nil {
				return err
			}
		}
	}

	if !onlyUsed && len(unused) > 0{
		fmt.Fprintln(writer, "\nUnused opcodes:")
		for _, str := range unused {
			fmt.Fprintln(writer, str)
		}
	}

	return nil
}
