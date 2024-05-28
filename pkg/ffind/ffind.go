// FFind implementation details.
package ffind

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/cuhsat/fact/internal/fact"
	"github.com/cuhsat/fact/internal/fact/hash"
	"github.com/cuhsat/fact/internal/fact/zip"
	"github.com/cuhsat/fact/internal/sys"
	"github.com/cuhsat/fact/pkg/windows"
)

const (
	fallback = "fact"
	extZip   = ".zip"
	extTxt   = ".txt"
	rLimit   = 1024
)

type ffind struct {
	wg sync.WaitGroup

	root string
	arc  string
	lst  string
	ha   string
	rp   bool
	so   bool
	uo   bool
}

type fstep func(in <-chan string, out chan<- string)

func Find(root, arc, lst, ha string, rp, so, uo bool) (files []string) {
	ff := &ffind{
		root: root,
		arc:  arc,
		lst:  lst,
		ha:   ha,
		rp:   rp,
		so:   so,
		uo:   uo,
	}

	// Switch to live mode
	if len(ff.root)+len(ff.arc)+len(ff.ha) == 0 {
		ff.live()
	}

	var ch [4]chan string
	var cn = 0

	for i := range ch {
		ch[i] = make(chan string, rLimit)
	}

	add := func(fn fstep) {
		ff.wg.Add(1)

		go fn(ch[cn], ch[cn+1])

		cn++
	}

	ff.wg.Add(1)

	go ff.enum(ch[cn])

	if len(ff.arc) > 0 {
		add(ff.comp)
	}

	if len(ff.lst) > 0 {
		add(ff.list)
	}

	if len(ff.ha) > 0 {
		add(ff.hash)
	}

	for l := range ch[cn] {
		if len(ff.ha) == 0 {
			l = ff.path(l)
		}

		files = append(files, l)
	}

	ff.wg.Wait()

	return
}

func (ff *ffind) enum(out chan<- string) {
	defer close(out)
	defer ff.wg.Done()

	if len(ff.root) > 0 {
		fi, err := os.Stat(ff.root)

		if err != nil {
			sys.Fatal(err)
		}

		if !fi.IsDir() {
			sys.Fatal("not a directory")
		}
	}

	if !ff.uo {
		windows.EnumSystem(ff.root, out)
	}

	if !ff.so {
		windows.EnumUsers(ff.root, out)
	}
}

func (ff *ffind) comp(in <-chan string, out chan<- string) {
	defer close(out)
	defer ff.wg.Done()

	z, err := zip.NewZip(ff.arc, time.Now().Format(time.RFC3339))

	if err != nil {
		sys.Fatal(err)
	}

	defer z.Close()

	for artifact := range in {
		err := z.Write(artifact, ff.path(artifact))

		if err != nil {
			sys.Error(err)
		}

		out <- artifact
	}
}

func (ff *ffind) list(in <-chan string, out chan<- string) {
	defer close(out)
	defer ff.wg.Done()

	f, err := os.Create(ff.lst)

	if err != nil {
		sys.Fatal(err)
	}

	defer f.Close()

	w := csv.NewWriter(f)

	err = w.Write(header())

	if err != nil {
		sys.Fatal(err)
	}

	for artifact := range in {
		err = w.Write(record(artifact))

		if err != nil {
			sys.Error(err)
		}

		out <- artifact
	}
}

func (ff *ffind) hash(in <-chan string, out chan<- string) {
	defer close(out)
	defer ff.wg.Done()

	for artifact := range in {
		s, err := hash.Sum(artifact, ff.ha)

		if err != nil {
			sys.Error(err)
			continue
		}

		out <- fmt.Sprintf("%x%s%s", s, fact.HashSep, ff.path(artifact))
	}
}

func (ff *ffind) live() {
	host, err := os.Hostname()

	if err != nil {
		sys.Error(err)

		host = fallback
	}

	ff.arc = host + extZip
	ff.lst = host + extTxt

	ff.rp = true
	ff.so = true
	ff.uo = true
}

func (ff *ffind) path(f string) string {
	if !ff.rp {
		return f
	}

	r := len(ff.root)

	if r > 0 {
		r++
	}

	return f[r:]
}

func header() []string {
	return []string{
		"Filename",
		"Path",
		"Size (bytes)",
		"Modified",
	}
}

func record(f string) (l []string) {
	l = append(l, filepath.Base(f))
	l = append(l, f)

	fi, err := os.Stat(f)

	if err != nil {
		sys.Error(err)
		return
	}

	l = append(l, fmt.Sprintf("%d", fi.Size()))
	l = append(l, fi.ModTime().String())

	return
}
