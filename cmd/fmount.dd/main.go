// Mount forensic raw or dd disk images for read-only processing.
//
// Usage:
//
//	fmount.dd [-fsuzqhv] [-H CRC32|MD5|SHA1|SHA256] [-V SUM] [-B KEY] [-D DIR] IMAGE
//
// The flags are:
//
//	 -D directory
//		The mount point directory.
//	 -B key
//	 	The BitLocker key.
//	 -H algorithm
//	 	The hash algorithm to use.
//	 -V sum
//	 	The hash sum to verify.
//	 -f
//		Force type (bypass check).
//	 -s
//		System partition only.
//	 -u
//		Unmount image.
//	 -z
//		Unzip image.
//	 -q
//		Quiet mode.
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

	"github.com/cuhsat/fact/internal/fact"
	"github.com/cuhsat/fact/internal/sys"
	"github.com/cuhsat/fact/pkg/fmount"
	"github.com/cuhsat/fact/pkg/fmount/dd"
)

func main() {
	D := flag.String("D", "", "Mount point")
	B := flag.String("B", "", "BitLocker key")
	H := flag.String("H", "", "Hash algorithm")
	V := flag.String("V", "", "Hash sum")
	f := flag.Bool("f", false, "Force mounting")
	s := flag.Bool("s", false, "System partition only")
	u := flag.Bool("u", false, "Unmount image")
	z := flag.Bool("z", false, "Unzip image")
	q := flag.Bool("q", false, "Quiet mode")
	h := flag.Bool("h", false, "Show usage")
	v := flag.Bool("v", false, "Show version")

	flag.CommandLine.SetOutput(io.Discard)
	flag.Parse()

	img := sys.Arg()

	if *v {
		sys.Final("fmount.dd", fact.Version)
	}

	if *h || len(img) == 0 {
		sys.Usage("fmount.dd [-fsuzqhv] [-H CRC32|MD5|SHA1|SHA256] [-V SUM] [-B KEY] [-D DIR] IMAGE")
	}

	if *q {
		sys.Progress = nil
	}

	if *z {
		ex, err := fmount.Extract(img)

		if err != nil {
			sys.Fatal(err)
		} else {
			img = ex
		}
	}

	if (len(*H) == 0) != (len(*V) == 0) {
		sys.Fatal("hash algorithm and sum are required")
	}

	if len(*H) > 0 && len(*V) > 0 {
		ok, err := fmount.Verify(img, *H, *V)

		if err != nil {
			sys.Fatal(err)
		}

		if !ok {
			sys.Fatal("hash sum does not match")
		}
	}

	if !*f {
		is, err := dd.Is(img)

		if err != nil {
			sys.Fatal(err)
		}

		if !is {
			sys.Fatal("image type not supported")
		}
	}

	var err error

	if *u {
		err = dd.Unmount(img)
	} else {
		_, err = dd.Mount(img, *D, *B, *s)
	}

	if err != nil {
		sys.Fatal(err)
	}
}
