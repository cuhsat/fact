// FMount functions.
package fmount

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuhsat/fact/internal/sys"
)

const (
	SymlinkPath = "/tmp/fmount"
)

func Dev(loop string) string {
	return filepath.Join("/dev", loop)
}

func BaseFile(name string) string {
	b := filepath.Base(name)

	return strings.TrimSuffix(b, filepath.Ext(b))
}

func Loops(img string) (loops []string, err error) {
	loops, err = LoSetupList(img)

	if err != nil {
		return
	}

	if len(loops) == 0 {
		err = errors.New("no devices found")
		return
	}

	return
}

func Parts(dev string) (parts []string, err error) {
	parts, err = LsBlk(dev, "name")

	if err != nil {
		return
	}

	if len(parts) <= 1 {
		err = errors.New("no partitions found")
		return
	}

	parts = parts[1:] // skip loop root

	return
}

func Mounts(dev string) (mnts []string, err error) {
	return LsBlk(dev, "mountpoints")
}

func IsEncrypted(dev string) (is bool, err error) {
	ft, err := LsBlk(dev, "fstype")

	if err != nil {
		return
	}

	if len(ft) != 1 {
		err = errors.New("type count differs")
		return
	}

	is = (ft[0] == "BitLocker")

	return
}

func CreateMount(img, mnt string) error {
	if len(mnt) == 0 {
		mnt = BaseFile(img)
	}

	return os.MkdirAll(mnt, sys.MODE_DIR)
}

func CreateSymlink(dev, mnt string) error {
	lnk := filepath.Join(SymlinkPath, dev)
	src := filepath.Join(mnt, DislockerDev)

	return os.Symlink(src, lnk)
}

func FollowSymlink(dev string) (src string, err error) {
	lnk := filepath.Join(SymlinkPath, dev)

	if _, err = os.Stat(lnk); os.IsNotExist(err) {
		err = errors.New("symlink not found")
		return
	}

	src, err = filepath.EvalSymlinks(lnk)

	return
}

func RemoveSymlink(dev string) error {
	lnk := filepath.Join(SymlinkPath, dev)

	return os.Remove(lnk)
}

func CreateDirf(root, m string, a ...any) (dir string, err error) {
	dir = filepath.Join(root, fmt.Sprintf(m, a...))
	err = os.MkdirAll(dir, sys.MODE_DIR)

	return
}

func RemoveDirs(dir string) (err error) {
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		return
	}

	ff, err := os.ReadDir(dir)

	if err != nil {
		return
	}

	for _, f := range ff {
		if !f.IsDir() {
			continue
		}

		os.Remove(filepath.Join(dir, f.Name()))
	}

	os.Remove(dir)

	return nil
}
