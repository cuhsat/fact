// LoSetup functions.
package fmount

import (
	"strings"

	"github.com/cuhsat/fact/internal/sys"
)

func LoSetupList(img string) (ls []string, err error) {
	out, err := sys.StdCall("losetup", "-l", "-n", "-O", "name", "-j", img)

	if err != nil {
		return
	}

	ls = strings.Split(strings.TrimSpace(out), "\n")

	return
}

func LoSetupAttach(img string) (dev string, err error) {
	dev, err = sys.StdCall("losetup", "-Pfr", "--show", img)

	return
}

func LoSetupDetach(dev string) (err error) {
	_, err = sys.StdCall("losetup", "-d", dev)

	return
}
