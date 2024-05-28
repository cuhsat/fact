// FMount implementation details.
package fmount

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuhsat/fact/internal/fact/hash"
	"github.com/cuhsat/fact/internal/fact/zip"
	"github.com/cuhsat/fact/internal/sys"
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
	i, err := zip.Index(img)

	if err != nil {
		return
	}

	if len(i) > 1 {
		err = errors.New("more than one file")
		return
	}

	dir := filepath.Dir(img)

	p = filepath.Join(dir, i[0])

	if _, err = os.Stat(p); !os.IsNotExist(err) {
		err = errors.New("file already exists")
		return
	}

	if err = zip.Unzip(img, dir); err != nil {
		return
	}

	return
}

func Verify(img, algo, sum string) (ok bool, err error) {
	b, err := hash.Sum(img, algo)

	if err != nil {
		return
	}

	ok = fmt.Sprintf("%x", b) == strings.ToLower(sum)

	return
}

func Forward(name string, arg ...string) {
	sys.ExitCall(name, arg...)
}
