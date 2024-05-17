// Mount forensic disk images for read-only processing.
//
// Usage:
//
//	fmount [-suzhv] [-T RAW|DD] [-D DIRECTORY] IMAGE
//
// The flags are:
//
//	 -D directory
//		The mount point directory.
//	 -T type
//	    The disk image type.
//	 -s
//		System partition only.
//	 -u
//		Unmount image.
//	 -z
//		Unzip image.
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

	"github.com/cuhsat/fact/internal/sys"
	"github.com/cuhsat/fact/pkg/fmount"
	"github.com/cuhsat/fact/pkg/fmount/dd"
)

// Changed by ldflags
var Version string = "dev"

func main() {
	D := flag.String("D", "", "Mount point")
	T := flag.String("T", "", "Image type")
	s := flag.Bool("s", false, "System partition only")
	u := flag.Bool("u", false, "Unmount image")
	z := flag.Bool("z", false, "Unzip image")
	h := flag.Bool("h", false, "Show usage")
	v := flag.Bool("v", false, "Show version")

	flag.CommandLine.SetOutput(io.Discard)
	flag.Parse()

	img := sys.Input()

	if *v {
		sys.Print("fmount", Version)
	}

	if *h || len(img) == 0 {
		sys.Usage("fmount [-suzhv] [-T RAW|DD] [-D DIRECTORY] IMAGE")
	}

	it, err := fmount.DetectType(img, *T)

	if err != nil {
		sys.Fatal(err)
	}

	args := make([]string, 0)

	if len(*D) > 0 {
		args = append(args, "-D", *D)
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

	args = append(args, img)

	switch it {
	case dd.RAW, dd.DD:
		fmount.Forward("fmount.dd", args...)
	default:
		sys.Fatal("image type not supported:", it)
	}
}
