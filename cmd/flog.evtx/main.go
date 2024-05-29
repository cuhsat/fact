// Log Windows event logs as JSON in ECS.
//
// Usage:
//
//	flog.evtx [-pqhv] [-D DIRECTORY] [FILE ...]
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

	"github.com/cuhsat/fact/internal/fact"
	"github.com/cuhsat/fact/internal/sys"
	"github.com/cuhsat/fact/pkg/flog"
	"github.com/cuhsat/fact/pkg/flog/evtx"
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
		sys.Final("flog.evtx", fact.Version)
	}

	if *h || len(files) == 0 {
		sys.Usage("flog.evtx [-pqhv] [-D DIRECTORY] [FILE ...]")
	}

	if *q {
		sys.Progress = nil
	}

	g := new(errgroup.Group)

	for _, f := range files {
		g.Go(func() (err error) {
			_, err = evtx.Log(f, *D, *p)
			return
		})
	}

	if err := g.Wait(); err != nil {
		sys.Fatal(err)
	}
}
