// Umount functions.
package fmount

import "github.com/cuhsat/fact/internal/sys"

func UmountDir(dir string) (err error) {
	_, err = sys.StdCall("umount", "-R", dir)

	return
}

func UmountDev(dev string) (err error) {
	_, err = sys.StdCall("umount", "-A", dev)

	return
}
