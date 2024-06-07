// QEMU functions.
package fmount

import (
	"github.com/cuhsat/fact/internal/sys"
)

const (
	QemuDev   = "/dev/nbd1"
	QemuParts = 6
)

func QemuAttach(dev, img string) (err error) {
	_, err = sys.StdCall("qemu-nbd", "--fork", "-r", "-c", dev, img)

	return
}

func QemuDetach(dev string) (err error) {
	_, err = sys.StdCall("qemu-nbd", "-d", dev)

	return
}
