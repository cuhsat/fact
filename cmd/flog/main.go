// Log forensic artifacts as JSON in ECS.
//
// Usage:
//
//	flog [-pqhv] [-D DIRECTORY] [FILE ...]
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
//		The artifact file(s) to process.
//		Defaults to STDIN if not given.
package main

import (
	"flag"
	"io"

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
		sys.Usage("flog [-pqhv] [-D DIRECTORY] [FILE ...]")
	}

	args := make([]string, 0)

	if len(*D) > 0 {
		args = append(args, "-D", *D)
	}

	if *p {
		args = append(args, "-p")
	}

	if *q {
		args = append(args, "-q")
	}

	g := new(errgroup.Group)

	g.Go(func() error {
		return flog.Evtx(files, args)
	})

	if err := g.Wait(); err != nil {
		sys.Fatal(err)
	}
}
