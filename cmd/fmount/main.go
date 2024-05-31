// Mount forensic disk images for read-only processing.
//
// Usage:
//
//	fmount [-suzqhv] [-H CRC32|MD5|SHA1|SHA256] [-V SUM] [-B KEY] [-T RAW|DD] [-D DIRECTORY] IMAGE
//
// The flags are:
//
//	 -D directory
//		The mount point directory.
//	 -T type
//	    The disk image type.
//	 -B key
//	 	The BitLocker key.
//	 -H algorithm
//	 	The hash algorithm to use.
//	 -V sum
//	 	The hash sum to verify.
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
	T := flag.String("T", "", "Image type")
	B := flag.String("B", "", "BitLocker key")
	H := flag.String("H", "", "Hash algorithm")
	V := flag.String("V", "", "Hash sum")
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
		sys.Final("fmount", fact.Version)
	}

	if *h || len(img) == 0 {
		sys.Usage("fmount [-suzqhv] [-H CRC32|MD5|SHA1|SHA256] [-V SUM] [-B KEY] [-T RAW|DD] [-D DIRECTORY] IMAGE")
	}

	it, err := fmount.DetectType(img, *T)

	if err != nil {
		sys.Fatal(err)
	}

	args := make([]string, 0)

	if len(*D) > 0 {
		args = append(args, "-D", *D)
	}

	if len(*B) > 0 {
		args = append(args, "-B", *B)
	}

	if len(*H) > 0 {
		args = append(args, "-H", *H)
	}

	if len(*V) > 0 {
		args = append(args, "-V", *V)
	}

	if *s {
		args = append(args, "-s")
	}

	if *u {
		args = append(args, "-u")
	}

	if *z {
		args = append(args, "-z")
	}

	if *q {
		args = append(args, "-q")
	}

	args = append(args, img)

	switch it {
	case dd.RAW, dd.DD:
		fmount.Forward("fmount.dd", args...)
	default:
		sys.Fatal("image type not supported:", it)
	}
}
