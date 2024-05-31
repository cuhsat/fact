// LsBlk functions.
package fmount

import (
	"errors"
	"strings"

	"github.com/cuhsat/fact/internal/sys"
)

func Parts(dev string) (parts []string, err error) {
	parts, err = LsBlk(dev, "name")

	if err != nil {
		return
	}

	if len(parts) <= 1 {
		err = errors.New("no partitions found")
		return
	}

	parts = parts[1:]

	return
}

func Mnts(dev string) (mnts []string, err error) {
	return LsBlk(dev, "mountpoints")
}

func LsBlk(dev, col string) (ls []string, err error) {
	out, err := sys.StdCall("lsblk", "-l", "-n", "-o", col, strings.TrimSpace(dev))

	if err != nil {
		return
	}

	ls = strings.Split(strings.TrimSpace(out), "\n")

	return
}
