// Log Windows event logs information in ECS schema.
//
// Usage:
//
//	flog.evtx [-hv] [-D DIRECTORY] [FILE ...]
//
// The flags are:
//
//	 -D directory
//	    The log directory.
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
	"strings"

	"github.com/cuhsat/fact/internal/fact"
	"github.com/cuhsat/fact/internal/sys"
	"github.com/cuhsat/fact/pkg/flog"
	"github.com/cuhsat/fact/pkg/flog/evtx"
	"golang.org/x/sync/errgroup"
)

func main() {
	D := flag.String("D", "", "Log directory")
	h := flag.Bool("h", false, "Show usage")
	v := flag.Bool("v", false, "Show version")

	flag.CommandLine.SetOutput(io.Discard)
	flag.Parse()

	files := flog.StripHash(sys.Args())

	if *v {
		sys.Print("flog.evtx", fact.Version)
	}

	if *h || len(files) == 0 {
		sys.Usage("flog.evtx [-hv] [-D DIRECTORY] [FILE ...]")
	}

	g := new(errgroup.Group)

	for _, f := range files {
		g.Go(func() (err error) {
			l, err := evtx.Log(f, *D)

			if err == nil {
				sys.Print(strings.Join(l, "\n"))
			}

			return
		})
	}

	if err := g.Wait(); err != nil {
		sys.Fatal(err)
	}
}
