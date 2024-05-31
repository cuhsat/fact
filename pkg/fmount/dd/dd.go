// DD implementation details.
package dd

import (
	"errors"
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
	if err = fmount.CreateMount(img, mnt); err != nil {
		return
	}

	// attach image as loop device
	loi, err := fmount.LoSetupAttach(img)

	if err != nil {
		return
	}

	// get partition loop devices
	lops, err := fmount.Parts(loi)

	if err != nil {
		return
	}

	// handle found partition loop devices
	for i, lop := range lops {
		dev := fmount.Dev(lop)

		// check if partition loop device is bootable
		sp, err := detectMagic(dev)

		if err != nil {
			sys.Error(err)
			continue
		}

		// if all or bootable
		if !so || (so && sp) {

			// create partition mount point
			mntp, err := fmount.CreateDirf(mnt, "p%d", i+1)

			if err != nil {
				sys.Error(err)
				continue
			}

			// check if partition loop device is encrypted
			is, err := fmount.IsEncrypted(dev)

			if err != nil {
				sys.Error(err)
				continue
			}

			// if encrypted
			if is {

				// check for key
				if len(key) == 0 {
					sys.Error("no key given")
					continue
				}

				// create fuse mount point
				mntf, err := fmount.CreateDirf(mnt, "p%d-fuse", i+1)

				if err != nil {
					sys.Error(err)
					continue
				}

				// mount to be decrypted partition loop device as fuse
				if err = fmount.DislockerFuse(dev, key, mntf); err != nil {
					sys.Error(err)
					continue
				}

				// create symlink to track device relations
				if err = fmount.CreateSymlink(lop, mntf); err != nil {
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
	lois, err := fmount.Loops(img)

	if err != nil {
		return
	}

	// handle found loop devices
	for _, loi := range lois {

		// get partition loop devices
		lops, err := fmount.Parts(loi)

		if err != nil {
			sys.Error(err)
			continue
		}

		// get mount points of partition loop device
		mnts, err := fmount.Mounts(loi)

		if err != nil {
			sys.Error(err)
			continue
		}

		// handle found loop devices
		for _, lop := range lops {
			dev := fmount.Dev(lop)

			// check if partition loop device is encrypted
			is, err := fmount.IsEncrypted(dev)

			if err != nil {
				sys.Error(err)
				continue
			}

			// if encrypted
			if is {

				// follow symlink and get partition mount points
				src, err := fmount.FollowSymlink(lop)

				if err != nil {
					sys.Error(err)
					continue
				}

				mntf := filepath.Dir(src)
				mntp := fuseRoot(src)
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
				if err = fmount.LoSetupDetach(lop); err != nil {
					sys.Error(err)
					continue
				}

				// remove symlink
				if err = fmount.RemoveSymlink(dev); err != nil {
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

func fuseRoot(dev string) string {
	return dev[:len(dev)-len("-fuse/"+fmount.DislockerDev)]
}
