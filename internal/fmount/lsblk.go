// LsBlk functions.
package fmount

import (
	"strings"

	"github.com/cuhsat/fact/internal/sys"
)

func LsBlk(dev, col string) (ls []string, err error) {
	out, err := sys.StdCall("lsblk", "-l", "-n", "-o", col, strings.TrimSpace(dev))

	if err != nil {
		return
	}

	ls = strings.Split(strings.TrimSpace(out), "\n")

	return
}
