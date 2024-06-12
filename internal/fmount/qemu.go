// QEMU functions.
package fmount

import (
	"github.com/cuhsat/fact/internal/sys"
)

const (
	QemuDev   = "/dev/nbd1"
	QemuParts = 16
)

func QemuAttach(dev, img string, xargs []string) (err error) {
	_, err = sys.StdCall("qemu-nbd", append([]string{
		"--fork", "-r", "-c", dev, img,
	}, xargs...)...)

	return
}

func QemuDetach(dev string) (err error) {
	_, err = sys.StdCall("qemu-nbd", "-d", dev)

	return
}
