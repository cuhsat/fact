// Find forensic artifacts in a mount point or on the live system.
//
// Usage:
//
//	ffind [-rsuqhv] [-H CRC32|MD5|SHA1|SHA256] [-Z ARCHIVE] [-F FILE] [MOUNT ...]
//
// The flags are:
//
//	 -H algorithm
//	 	The hash algorithm to use.
//	 -Z archive
//		The artifacts archive name.
//	 -F file
//		The filename to write also.
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
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cuhsat/fact/internal/sys"
	"github.com/cuhsat/fact/pkg/ffind"
)

// Changed by ldflags
var Version string = "dev"

func main() {
	H := flag.String("H", "", "Hash algorithm")
	Z := flag.String("Z", "", "Archive name")
	F := flag.String("F", "", "File to write")
	r := flag.Bool("r", false, "Relative paths")
	s := flag.Bool("s", false, "System artifacts only")
	u := flag.Bool("u", false, "User artifacts only")
	q := flag.Bool("q", false, "Quiet mode")
	h := flag.Bool("h", false, "Show usage")
	v := flag.Bool("v", false, "Show version")

	flag.CommandLine.SetOutput(io.Discard)
	flag.Parse()

	mnt := sys.Args()

	if *h {
		sys.Usage("ffind [-rsuqhv] [-H CRC32|MD5|SHA1|SHA256] [-Z ARCHIVE] [-F FILE] [MOUNT ...]")
	}

	if *v {
		sys.Print("ffind", Version)
	}

	if *q && len(*F)+len(*Z) == 0 {
		sys.Fatal("archive or file required")
	}

	if *r && len(mnt) > 1 {
		sys.Error("relative paths disabled")
		*r = false
	}

	for _, p := range mnt {
		files := ffind.Find(p, *Z, *H, *r, *s, *u)

		if len(*F) > 0 {
			f, err := os.OpenFile(*F, os.O_WRONLY|os.O_CREATE, 0666)

			if err != nil {
				sys.Error(err)
			}

			b := []byte(strings.Join(files, "\n"))

			if _, err = f.Write(b); err != nil {
				sys.Error(err)
			}

			if err = f.Close(); err != nil {
				sys.Error(err)
			}
		}

		if !*q {
			for _, f := range files {
				fmt.Fprintln(os.Stdout, f)
			}
		}
	}
}
