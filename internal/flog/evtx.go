// Evtx functions.
package flog

import (
	"os"
	"path/filepath"

	"github.com/cuhsat/fact/internal/fact/ez"
	"github.com/cuhsat/fact/internal/sys"
)

func ImportEvent(src, dir string) (lines []string, err error) {
	log, err := evtxecmd(src, src+".json", dir)

	if err != nil {
		return
	}

	lines, err = ReadLines(log)

	if err != nil {
		return
	}

	err = os.Remove(log)

	return
}

func ExportEvent(b []byte, dst string) (err error) {
	f, err := os.Create(dst)

	if err != nil {
		return
	}

	_, err = f.Write(b)

	f.Close()

	return
}

func evtxecmd(src, dst, dir string) (log string, err error) {
	asm, err := ez.Path("EvtxECmd.dll")

	if err != nil {
		return
	}

	if len(dir) == 0 {
		dir = filepath.Dir(dst)
	}

	dst = filepath.Base(dst)
	log = filepath.Join(dir, dst)

	_, err = sys.StdCall("dotnet", asm, "-f", src, "--fj", "--json", dir, "--jsonf", dst)

	return
}
