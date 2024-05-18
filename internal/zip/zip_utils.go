// Zip utils functions.
package zip

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"path/filepath"
	"slices"
)

func Index(name string) (files []string, err error) {
	a, err := zip.OpenReader(name)

	if err != nil {
		return
	}

	defer a.Close()

	for _, f := range a.File {
		if !f.FileInfo().IsDir() {
			files = append(files, filepath.ToSlash(f.Name))
		}
	}

	slices.Sort(files)

	return
}

func Unzip(name, dir string) (err error) {
	a, err := zip.OpenReader(name)

	if err != nil {
		return
	}

	defer a.Close()

	for _, f := range a.File {
		file := path.Join(dir, filepath.ToSlash(f.Name))

		if f.FileInfo().IsDir() {
			os.MkdirAll(file, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(path.Dir(file), os.ModePerm); err != nil {
			return err
		}

		src, err := f.Open()

		if err != nil {
			return err
		}

		dst, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())

		if err != nil {
			src.Close()
			return err
		}

		_, err = io.Copy(dst, src)

		dst.Close()
		src.Close()

		if err != nil {
			return err
		}
	}

	return
}
