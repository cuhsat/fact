// File functions.
package flog

import (
	"bufio"
	"os"
)

func ReadLines(name string) (lines []string, err error) {
	f, err := os.Open(name)

	if err != nil {
		return
	}

	defer f.Close()

	fs := bufio.NewScanner(f)

	fs.Split(bufio.ScanLines)

	for fs.Scan() {
		lines = append(lines, fs.Text())
	}

	return
}
