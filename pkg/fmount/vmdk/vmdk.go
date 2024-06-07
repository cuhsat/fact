// VMDK implementation details.
package vmdk

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cuhsat/fact/internal/fmount"
	"github.com/cuhsat/fact/internal/sys"
)

const (
	VMDK = "vmdk"
)

var (
	header = []byte{0x4B, 0x44, 0x4D, 0x56}
)

func Is(img string) (is bool, err error) {
	f, err := os.Open(img)

	if err != nil {
		return
	}

	defer f.Close()

	b := make([]byte, 4)

	n, err := f.Read(b)

	if err != nil {
		return
	}

	if n != len(b) {
		err = errors.New("could not read header")
		return
	}

	return bytes.Equal(b, header), nil
}

func Mount(img, mnt, key string, so bool) (parts []string, err error) {

	// create symlink directory
	if err = os.MkdirAll(fmount.SymlinkPath, sys.MODE_DIR); err != nil {
		return
	}

	// create mount point
	if err = fmount.CreateImageMount(img, mnt); err != nil {
		return
	}

	// ensure kernel module is loaded
	if err = ensureMod(fmount.QemuParts); err != nil {
		return
	}

	// attach image as network block device
	if err = fmount.QemuAttach(fmount.QemuDev, img); err != nil {
		return
	}

	// create image symlink directory to track image relations
	if err = fmount.CreateImageSymlink(img, fmount.QemuDev); err != nil {
		return
	}

	// get partition network block devices
	nbdps, err := fmount.PartDevs(fmount.QemuDev)

	if err != nil {
		return
	}

	// handle found partitions
	for i, nbdp := range nbdps {
		dev := fmount.Dev(nbdp)

		// check if partition is bootable
		sp, err := fmount.IsBootable(dev)

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

			// check if partition is encrypted
			is, err := fmount.IsEncrypted(dev)

			if err != nil {
				sys.Error(err)
				continue
			}

			// if encrypted
			if is {

				// check if key given
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

				// mount to be decrypted partition as fuse
				if err = fmount.DislockerFuse(dev, key, mntf); err != nil {
					sys.Error(err)
					continue
				}

				// create symlink to track device relations
				if err = fmount.CreateSymlink(nbdp, mntf); err != nil {
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

	// get network block devices associated with image
	nbds, err := fmount.BlockDevs(img)

	if err != nil {
		return
	}

	// handle found network block devices
	for _, nbd := range nbds {

		// get partition devices
		nbdps, err := fmount.PartDevs(nbd)

		if err != nil {
			sys.Error(err)
			continue
		}

		// get mount points of device
		mnts, err := fmount.Mounts(nbd)

		if err != nil {
			sys.Error(err)
			continue
		}

		// handle found partitions
		for _, nbdp := range nbdps {
			dev := fmount.Dev(nbdp)

			// check if partition is encrypted
			is, err := fmount.IsEncrypted(dev)

			if err != nil {
				sys.Error(err)
				continue
			}

			// if encrypted
			if is {

				// follow symlink and get partition mount points
				src, err := fmount.FollowSymlink(nbdp)

				if err != nil {
					sys.Error(err)
					continue
				}

				mntf := filepath.Dir(src)
				mntp := fmount.FromFuse(src)
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

				// detach partition network block device
				if err = fmount.QemuDetach(dev); err != nil {
					sys.Error(err)
					continue
				}

				// remove symlink
				if err = fmount.RemoveSymlink(nbdp); err != nil {
					sys.Error(err)
					continue
				}
			} else {

				// unmount partition network block device
				if err = fmount.UmountDev(dev); err != nil {
					sys.Error(err)
					continue
				}

				// detach partition network block device
				if err = fmount.QemuDetach(dev); err != nil {
					sys.Error(err)
					continue
				}
			}
		}

		// remove image symlink directory
		dir := filepath.Join(fmount.SymlinkPath, filepath.Base(img))

		if err = os.RemoveAll(dir); err != nil {
			sys.Error(err)
		}

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

func ensureMod(mp int) (err error) {
	is, err := fmount.IsLoaded("nbd")

	if err != nil {
		return
	}

	if !is {
		fmount.ModLoad("nbd", fmt.Sprintf("max_part=%d", mp))
	}

	return
}
