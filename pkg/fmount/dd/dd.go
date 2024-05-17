// DD implementation details.
package dd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cuhsat/fact/internal/sys"
)

const (
	DD   = "dd"
	RAW  = "raw"
	mode = 0444
)

func Is(img string) (is bool, err error) {
	return detectMagic(img)
}

func Mount(img, dir string, so bool) (sysroot string, err error) {
	if len(dir) == 0 {
		b := filepath.Base(img)

		dir = strings.TrimSuffix(b, filepath.Ext(b))
	}

	if err = os.MkdirAll(dir, mode); err != nil {
		return
	}

	lo, err := losetupAttach(img)

	if err != nil {
		return
	}

	ls, err := lsblk(lo, "name")

	if err != nil {
		return
	}

	if len(ls) <= 1 {
		err = errors.New("no partitions found")
		return
	}

	for i, l := range ls[1:] {
		p := path.Join(dir, fmt.Sprintf("p%d", i+1))

		if err := os.MkdirAll(p, mode); err != nil {
			sys.Error(err)
			continue
		}

		d := path.Join("/dev", l)

		sp, err := detectMagic(d)

		if err != nil {
			sys.Error(err)
			continue
		}

		if !so || (so && sp) {
			if mount(d, p) != nil {
				sys.Error(err)
				continue
			}
		}

		if sp {
			sysroot, err = filepath.Abs(p)

			if err != nil {
				sys.Error(err)
			}
		}
	}

	return
}

func Unmount(img string) (err error) {
	img, err = filepath.Abs(img)

	if err != nil {
		return
	}

	lo, err := losetupList(img)

	if err != nil {
		return
	}

	if len(lo) == 0 {
		return errors.New("no devices found")
	}

	for _, l := range lo {
		ls, err := lsblk(l, "name")

		if err != nil {
			sys.Error(err)
			continue
		}

		if len(ls) <= 1 {
			return errors.New("no partitions found")
		}

		mp, err := lsblk(l, "mountpoints")

		if err != nil {
			sys.Error(err)
			continue
		}

		for _, d := range ls[1:] {
			if err = umount(path.Join("/dev", d)); err != nil {
				sys.Error(err)
				continue
			}
		}

		if losetupDetach(l) != nil {
			sys.Error(err)
			continue
		}

		for _, p := range mp {
			if _, err := os.Stat(p); !os.IsNotExist(err) {
				if err = os.Remove(p); err != nil {
					sys.Error(err)
					continue
				}

				if err = os.Remove(filepath.Dir(p)); err != nil {
					sys.Error(err)
				}
			}
		}
	}

	return
}

func detectMagic(name string) (has bool, err error) {
	f, err := os.Open(name)

	if err != nil {
		return
	}

	defer f.Close()

	s := make([]byte, 512)

	n, err := f.Read(s)

	if err != nil {
		return
	}

	if n != len(s) {
		err = errors.New("could not read sector")
		return
	}

	return s[0x1FE] == 0x55 && s[0x1FF] == 0xAA, nil
}

func losetupAttach(img string) (dev string, err error) {
	dev, err = sys.StdCall("losetup", "-Pfr", "--show", img)

	return
}

func losetupDetach(dev string) (err error) {
	_, err = sys.StdCall("losetup", "-d", dev)

	return
}

func losetupList(img string) (l []string, err error) {
	lo, err := sys.StdCall("losetup", "-l", "-n", "-O", "name", "-j", img)

	if err != nil {
		return
	}

	l = strings.Split(strings.TrimSpace(lo), "\n")

	return
}

func lsblk(dev, col string) (l []string, err error) {
	ls, err := sys.StdCall("lsblk", "-l", "-n", "-o", col, strings.TrimSpace(dev))

	if err != nil {
		return
	}

	l = strings.Split(strings.TrimSpace(ls), "\n")

	return
}

func mount(dev, dir string) (err error) {
	_, err = sys.StdCall("mount", "-o", "ro", dev, dir)

	return
}

func umount(dev string) (err error) {
	_, err = sys.StdCall("umount", "-A", dev)

	return
}
