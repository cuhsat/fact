// Find forensic artifacts in mount points or on the live system.
//
// Usage:
//
//	ffind [-rsuqhv] [-H CRC32|MD5|SHA1|SHA256] [-Z ARCHIVE] [-L FILE] [MOUNT ...]
//
// The flags are:
//
//	 -H algorithm
//	 	The hash algorithm to use.
//	 -Z archive
//		The artifacts archive name.
//	 -L file
//		The artifacts listing name.
//	 -r
//		Output relative paths.
//	 -s
//		System artifacts only.
//	 -u
//		User artifacts only.
//	 -q
//		Quiet mode.
//	 -h
//		Show usage.
//	 -v
//		Show version.
//
// The arguments are:
//
//	 mount
//		The image mount point(s) or the system root path(s).
//		Defaults to STDIN, then %SYSTEMDRIVE% if not given.
package main

import (
	"flag"
	"io"

	"github.com/cuhsat/fact/internal/fact"
	"github.com/cuhsat/fact/internal/sys"
	"github.com/cuhsat/fact/pkg/ffind"
)

func main() {
	H := flag.String("H", "", "Hash algorithm")
	Z := flag.String("Z", "", "Archive name")
	L := flag.String("L", "", "Listing name")
	r := flag.Bool("r", false, "Relative paths")
	s := flag.Bool("s", false, "System artifacts only")
	u := flag.Bool("u", false, "User artifacts only")
	q := flag.Bool("q", false, "Quiet mode")
	h := flag.Bool("h", false, "Show usage")
	v := flag.Bool("v", false, "Show version")

	flag.CommandLine.SetOutput(io.Discard)
	flag.Parse()

	mnt := sys.Args()

	if *v {
		sys.Final("ffind", fact.Version)
	}

	if *h {
		sys.Usage("ffind [-rsuqhv] [-H CRC32|MD5|SHA1|SHA256] [-Z ARCHIVE] [-L FILE] [MOUNT ...]")
	}

	if *q {
		sys.Progress = nil
	}

	if *q && len(*Z)+len(*L) == 0 {
		sys.Fatal("archive or listing required")
	}

	if *r && len(mnt) > 1 {
		sys.Error("relative paths disabled")
		*r = false
	}

	for _, p := range mnt {
		ffind.Find(p, *Z, *L, *H, *r, *s, *u)
	}
}
