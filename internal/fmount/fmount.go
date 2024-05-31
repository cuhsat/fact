// FMount functions.
package fmount

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	SymlinkPath = "/tmp/fmount"
)

func BaseFile(name string) string {
	b := filepath.Base(name)

	return strings.TrimSuffix(b, filepath.Ext(b))
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
