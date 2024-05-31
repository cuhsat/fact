// DD implementation details.
package dd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cuhsat/fact/internal/fmount"
	"github.com/cuhsat/fact/internal/sys"
)

const (
	DD  = "dd"
	RAW = "raw"
)

func Is(img string) (is bool, err error) {
	return detectMagic(img)
}

func Mount(img, mnt, key string, so bool) (parts []string, err error) {

	// create symlink directory
	if err = os.MkdirAll(fmount.SymlinkPath, sys.MODE_DIR); err != nil {
		return
	}

	// create mount point
	if len(mnt) == 0 {
		mnt = fmount.BaseFile(img)
	}

	if err = os.MkdirAll(mnt, sys.MODE_DIR); err != nil {
		return
	}

	// attach image as loop device
	loi, err := fmount.LoSetupAttach(img)

	if err != nil {
		return
	}

	// get partition loop devices
	lops, err := fmount.LsBlk(loi, "name")

	if err != nil {
		return
	}

	if len(lops) <= 1 {
		err = errors.New("no partitions found")
		return
	}

	// handle found partition loop devices (but skipping root)
	for i, lop := range lops[1:] {
		dev := toDev(lop)

		// check if partition loop device is bootable
		sp, err := detectMagic(dev)

		if err != nil {
			sys.Error(err)
			continue
		}

		// if all or bootable
		if !so || (so && sp) {

			// create partition mount point
			mntp := filepath.Join(mnt, fmt.Sprintf("p%d", i+1))

			if err := os.MkdirAll(mntp, sys.MODE_DIR); err != nil {
				sys.Error(err)
				continue
			}

			// check if partition loop device is encrypted
			is, err := fmount.IsEncrypted(dev)

			if err != nil {
				sys.Error(err)
				continue
			}

			if is && len(key) == 0 {
				sys.Error("no key given")
				continue
			}

			// if encrypted
			if is {

				// create fuse mount point
				mntf := filepath.Join(mnt, fmt.Sprintf("p%d-fuse", i+1))

				if err = os.MkdirAll(mntf, sys.MODE_DIR); err != nil {
					sys.Error(err)
					continue
				}

				// mount to be decrypted partition loop device as fuse
				err := fmount.DislockerFuse(dev, key, mntf)

				if err != nil {
					sys.Error(err)
					continue
				}

				// create symlink to track device relations
				src := filepath.Join(mntf, fmount.DislockerDev)
				lnk := filepath.Join(fmount.SymlinkPath, lop)

				if err = os.Symlink(src, lnk); err != nil {
					sys.Error(err)
					continue
				}

				// overwrite device to be mounted
				dev = filepath.Join(mntf, fmount.DislockerDev)
			}

			// mount device
			if fmount.Mount(dev, mntp, is) != nil {
				sys.Error(err)
				continue
			}

			// report progress
			mntp, err = filepath.Abs(mntp)

			if err != nil {
				sys.Error(err)
				continue
			}

			if sys.Progress != nil {
				sys.Progress(mntp)
			}

			parts = append(parts, mntp)
		}
	}

	return parts, nil
}

func Unmount(img string) (err error) {
	img, err = filepath.Abs(img)

	if err != nil {
		return
	}

	// get loop devices associated with image
	lois, err := fmount.LoSetupList(img)

	if err != nil {
		return
	}

	if len(lois) == 0 {
		return errors.New("no devices found")
	}

	// handle found loop devices
	for _, loi := range lois {

		// get partition loop devices
		lops, err := fmount.LsBlk(loi, "name")

		if err != nil {
			sys.Error(err)
			continue
		}

		if len(lops) <= 1 {
			sys.Error("no partitions found")
			continue
		}

		// get mount points of partition loop device
		mnts, err := fmount.LsBlk(loi, "mountpoints")

		if err != nil {
			sys.Error(err)
			continue
		}

		// handle found loop devices (but skipping root)
		for _, lop := range lops[1:] {
			dev := toDev(lop)

			// check if partition loop device is encrypted
			is, err := fmount.IsEncrypted(dev)

			if err != nil {
				sys.Error(err)
				continue
			}

			// if encrypted
			if is {

				// get symlink to fuse device for partition loop device
				lnk := filepath.Join(fmount.SymlinkPath, lop)

				if _, err := os.Stat(lnk); os.IsNotExist(err) {
					sys.Error("fmount symlink not found")
					continue
				}

				// get partition mount points from symlink
				src, err := filepath.EvalSymlinks(lnk)

				if err != nil {
					sys.Error(err)
					continue
				}

				mntf := filepath.Dir(src)
				mntp := toDir(src)

				mnts = append(mnts, mntp)

				// unmount partition mount point
				if err = fmount.UmountDir(mntp); err != nil {
					sys.Error(err)
					continue
				}

				// unmount fuse mount point
				if err = fmount.UmountDir(mntf); err != nil {
					sys.Error(err)
					continue
				}

				// detach partitionloop device
				if fmount.LoSetupDetach(lop) != nil {
					sys.Error(err)
					continue
				}

				// remove symlink
				if err = os.Remove(lnk); err != nil {
					sys.Error(err)
					continue
				}
			} else {

				// unmount partition loop device
				if err = fmount.UmountDev(dev); err != nil {
					sys.Error(err)
					continue
				}
			}
		}

		// detach loop device
		fmount.LoSetupDetach(loi)

		// remove empty mount points
		for _, mnt := range mnts {
			err = fmount.RemoveDirs(filepath.Dir(mnt))

			if err != nil {
				sys.Error(err)
				continue
			}
		}
	}

	return nil
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

func toDev(lo string) (dev string) {
	return filepath.Join("/dev", lo)
}

func toDir(lnk string) (dir string) {
	return lnk[:len(lnk)-len("-fuse/"+fmount.DislockerDev)]
}
