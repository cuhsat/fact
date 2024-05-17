// FMount implementation details.
package fmount

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cuhsat/fact/internal/sys"
	"github.com/cuhsat/fact/internal/zip"
	"github.com/cuhsat/fact/pkg/fmount/dd"
)

func DetectType(img, t string) (string, error) {
	e := filepath.Ext(img)

	if len(t) > 0 {
		return strings.ToLower(t), nil
	} else if len(e) > 0 {
		return strings.ToLower(e[1:]), nil
	} else if is, err := dd.Is(img); is {
		return dd.RAW, err
	} else {
		return "unknown", err
	}
}

func Extract(img string) (p string, err error) {
	dir, err := os.Getwd()

	if err != nil {
		return
	}

	i, err := zip.Index(img)

	if err != nil {
		return
	}

	if len(i) > 1 {
		err = errors.New("more than one file")
		return
	}

	if err = zip.Unzip(img, dir); err != nil {
		return
	}

	p = path.Join(dir, i[0])

	return
}

func Forward(name string, arg ...string) {
	sys.ExitCall(name, arg...)
}
