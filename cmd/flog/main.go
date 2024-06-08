// Log forensic artifacts as JSON in ECS.
//
// Usage:
//
//	flog [-pqhv] [-D DIR] [FILE ...]
//
// The flags are:
//
//	 -D directory
//	    The log directory.
//	 -p
//		Pretty JSON.
//	 -q
//		Quiet mode.
//	 -h
//		Show usage.
//	 -v
//		Show version.
//
// The arguments are:
//
//	 file
//		The event log file(s) to process.
//		Defaults to STDIN if not given.
package main

import (
	"flag"
	"io"
	"path/filepath"

	"github.com/cuhsat/fact/internal/fact"
	"github.com/cuhsat/fact/internal/sys"
	"github.com/cuhsat/fact/pkg/flog"
	"golang.org/x/sync/errgroup"
)

func main() {
	D := flag.String("D", "", "Log directory")
	p := flag.Bool("p", false, "Pretty JSON")
	q := flag.Bool("q", false, "Quiet mode")
	h := flag.Bool("h", false, "Show usage")
	v := flag.Bool("v", false, "Show version")

	flag.CommandLine.SetOutput(io.Discard)
	flag.Parse()

	files := flog.StripHash(sys.Args())

	if *v {
		sys.Final("flog", fact.Version)
	}

	if *h || len(files) == 0 {
		sys.Usage("flog [-pqhv] [-D DIR] [FILE ...]")
	}

	if *q {
		sys.Progress = nil
	}

	g := new(errgroup.Group)

	for _, f := range files {
		if filepath.Ext(f) == flog.Evtx {
			g.Go(func() (err error) {
				_, err = flog.LogEvent(f, *D, *p)
				return
			})
		}
	}

	if err := g.Wait(); err != nil {
		sys.Fatal(err)
	}
}
