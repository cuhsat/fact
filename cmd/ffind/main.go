// Find forensic artifacts in a mount point or on the live system.
//
// Usage:
//
//	ffind [-rsuqhv] [-H CRC32|MD5|SHA1|SHA256] [-Z ARCHIVE] [-F FILE] [SYSROOT]
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
//	 sysroot
//		The systems root path or image mount point.
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

	if *h {
		sys.Usage("ffind [-rsuqhv] [-H CRC32|MD5|SHA1|SHA256] [-Z ARCHIVE] [-F FILE] [SYSROOT]")
	}

	if *v {
		sys.Print("ffind", Version)
	}

	if *q && len(*Z) == 0 {
		sys.Fatal("archive name required")
	}

	files := ffind.Find(sys.Param(), *Z, *H, *r, *s, *u)

	if len(*F) > 0 {
		b := []byte(strings.Join(files, "\n"))

		if err := os.WriteFile(*F, b, 0666); err != nil {
			sys.Error(err)
		}
	}

	if !*q {
		for _, f := range files {
			fmt.Fprintln(os.Stdout, f)
		}
	}
}
