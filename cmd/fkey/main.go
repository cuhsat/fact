// Shows all BitLocker Recovery Key IDs of an image.
//
// Usage:
//
//	fkey [-hv] IMAGE
//
// The flags are:
//
//	 -h
//		Show usage.
//	 -v
//		Show version.
//
// The arguments are:
//
//	 image
//		The disk images filename.
package main

import (
	"flag"
	"io"
	"strings"

	"github.com/cuhsat/fact/internal/fact"
	"github.com/cuhsat/fact/internal/sys"
	"github.com/cuhsat/fact/pkg/fkey"
)

func main() {
	h := flag.Bool("h", false, "Show usage")
	v := flag.Bool("v", false, "Show version")

	flag.CommandLine.SetOutput(io.Discard)
	flag.Parse()

	img := sys.Arg()

	if *v {
		sys.Final("fkey", fact.Version)
	}

	if *h || len(img) == 0 {
		sys.Usage("fkey [-hv] IMAGE")
	}

	ids, err := fkey.RecoveryKeyIds(img)

	if err != nil {
		sys.Fatal(err)
	} else {
		sys.Final(strings.Join(ids, "\n"))
	}
}
