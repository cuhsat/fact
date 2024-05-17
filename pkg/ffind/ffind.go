// FFind implementation details.
package ffind

import (
	"archive/zip"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cuhsat/fact/internal/sys"
	"github.com/cuhsat/fact/pkg/windows"
)

const (
	CRC32  = "crc32"
	MD5    = "md5"
	SHA1   = "sha1"
	SHA256 = "sha256"
	SHA512 = "sha512"
	rLimit = 1024
)

type ffind struct {
	wg sync.WaitGroup

	sysroot string
	archive string
	hash    string
	rp      bool
	so      bool
	uo      bool
}

func Find(sysroot, archive, hash string, rp, so, uo bool) (lines []string) {
	// Going into live mode
	if len(sysroot)+len(archive)+len(hash) == 0 {
		host, err := os.Hostname()

		if err != nil {
			sys.Error(err)
			host = "fact" // fallback
		}

		archive = host + ".zip"
		hash = SHA256
	}

	ff := &ffind{
		sysroot: sysroot,
		archive: archive,
		hash:    hash,
		rp:      rp,
		so:      so,
		uo:      uo,
	}

	ch1 := make(chan string, rLimit)
	ch2 := make(chan string, rLimit)
	ch3 := make(chan string, rLimit)

	if len(ff.archive) > 0 {
		ff.wg.Add(3)

		go ff.find(ch1)
		go ff.zip(ch1, ch2)
		go ff.log(ch2, ch3)
	} else {
		ff.wg.Add(2)

		go ff.find(ch1)
		go ff.log(ch1, ch3)
	}

	for l := range ch3 {
		lines = append(lines, l)
	}

	ff.wg.Wait()

	return
}

func (ff *ffind) find(out chan<- string) {
	defer close(out)
	defer ff.wg.Done()

	if len(ff.sysroot) > 0 {
		fi, err := os.Stat(ff.sysroot)

		if err != nil {
			sys.Fatal(err)
		}

		if !fi.IsDir() {
			sys.Fatal("not a directory")
		}
	}

	if !ff.uo {
		windows.EnumSystem(ff.sysroot, out)
	}

	if !ff.so {
		windows.EnumUsers(ff.sysroot, out)
	}
}

func (ff *ffind) zip(in, out chan string) {
	defer close(out)
	defer ff.wg.Done()

	// TODO: file init after something was found
	a, err := os.Create(ff.archive)

	if err != nil {
		sys.Error(err)
		return
	}

	defer a.Close()

	w := zip.NewWriter(a)

	defer w.Close()

	w.SetComment(time.Now().Format(time.RFC3339))

	for artifact := range in {

		src, err := os.Open(artifact)

		if err != nil {
			sys.Error(err)
			continue
		}

		dst, err := w.Create(ff.path(artifact))

		if err != nil {
			sys.Error(err)
			src.Close()
			continue
		}

		_, err = io.Copy(dst, src)

		src.Close()

		if err != nil {
			sys.Error(err)
			continue
		}

		out <- artifact
	}
}

func (ff *ffind) log(in, out chan string) {
	defer close(out)
	defer ff.wg.Done()

	var h hash.Hash

	switch strings.ToLower(ff.hash) {
	case CRC32:
		h = crc32.NewIEEE()
	case MD5:
		h = md5.New()
	case SHA1:
		h = sha1.New()
	case SHA256:
		h = sha256.New()
	case SHA512:
		h = sha512.New()
	default:
		if len(ff.hash) > 0 {
			sys.Error("hash not supported:", ff.hash)
		}
	}

	if h == nil {
		for artifact := range in {
			out <- ff.path(artifact)
		}

		return
	}

	for artifact := range in {
		h.Reset()

		src, err := os.Open(artifact)

		if err != nil {
			sys.Error(err)
			continue
		}

		_, err = io.Copy(h, src)

		src.Close()

		if err != nil {
			sys.Error(err)
			continue
		}

		out <- fmt.Sprintf("%x %s", h.Sum(nil), ff.path(artifact))
	}
}

func (ff *ffind) path(f string) string {
	if !ff.rp {
		return f
	}

	r := len(ff.sysroot)

	if r > 0 {
		r++
	}

	return f[r:]
}
