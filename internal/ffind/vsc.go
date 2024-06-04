//go:build windows

// Volume Shadow Copy functions (Windows only).
package ffind

import (
	"os"
	"path/filepath"

	"github.com/mxk/go-vss"
)

const (
	symlink = "ffind"
)

func ShadowCopy(drv string) (dir string, err error) {
	dir = filepath.Join(os.TempDir(), symlink)

	if _, err = os.Stat(dir); !os.IsNotExist(err) {
		os.Remove(dir)
	}

	err = vss.CreateLink(dir, drv)

	return
}
