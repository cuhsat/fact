// FFind implementation details.
package ffind

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/cuhsat/fact/internal/hash"
	"github.com/cuhsat/fact/internal/sys"
	"github.com/cuhsat/fact/internal/zip"
	"github.com/cuhsat/fact/pkg/windows"
)

const (
	liveZip = "fact"
	liveExt = ".zip"
	rLimit  = 1024
)

type ffind struct {
	wg sync.WaitGroup

	sysroot string
	archive string
	algo    string
	rp      bool
	so      bool
	uo      bool
}

func Find(sysroot, archive, algo string, rp, so, uo bool) (lines []string) {
	ff := &ffind{
		sysroot: sysroot,
		archive: archive,
		algo:    algo,
		rp:      rp,
		so:      so,
		uo:      uo,
	}

	// Go into live mode
	if len(ff.sysroot)+len(ff.archive)+len(ff.algo) == 0 {
		ff.live()
	}

	ch1 := make(chan string, rLimit)
	ch2 := make(chan string, rLimit)
	ch3 := make(chan string, rLimit)

	ff.wg.Add(3)

	go ff.find(ch1)
	go ff.zip(ch1, ch2)
	go ff.log(ch2, ch3)

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

func (ff *ffind) zip(in <-chan string, out chan<- string) {
	defer close(out)
	defer ff.wg.Done()

	var z *zip.Zip
	var err error

	for artifact := range in {
		if len(ff.archive) > 0 {
			if z == nil { // init once
				meta := time.Now().Format(time.RFC3339)

				z, err = zip.NewZip(ff.archive, meta)

				if err != nil {
					sys.Fatal(err)
				}

				defer z.Close()
			}

			err := z.Write(artifact, ff.path(artifact))

			if err != nil {
				sys.Error(err)
			}
		}

		out <- artifact
	}
}

func (ff *ffind) log(in <-chan string, out chan<- string) {
	defer close(out)
	defer ff.wg.Done()

	for artifact := range in {
		p := ff.path(artifact)

		if len(ff.algo) > 0 {
			s, err := hash.Sum(artifact, ff.algo)

			if err != nil {
				sys.Error(err)
				continue
			}

			out <- fmt.Sprintf("%x  %s", s, p)
		} else {
			out <- p
		}
	}
}

func (ff *ffind) live() {
	host, err := os.Hostname()

	if err != nil {
		sys.Error(err)

		host = liveZip // fallback
	}

	ff.archive = host + liveExt
	ff.rp = true
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
