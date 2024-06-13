// Eric Zimmermann tools.
package flog

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cuhsat/fact/internal/fact/ez"
	"github.com/cuhsat/fact/internal/sys"
)

func Evtxe(src, dir string) (log string, err error) {
	cmd, err := ez.Path("EvtxECmd.dll")

	if err != nil {
		return
	}

	if len(dir) == 0 {
		dir = filepath.Dir(src)
	}

	dst := filepath.Base(src) + ".json"
	log = filepath.Join(dir, dst)

	_, err = sys.StdCall("dotnet", cmd, "-f", src, "--fj", "--json", dir, "--jsonf", dst)

	return
}

func Jle(src, dir string) (log string, err error) {
	cmd, err := ez.Path("JLECmd.dll")

	if err != nil {
		return
	}

	if len(dir) == 0 {
		dir = filepath.Dir(src)
	}

	dst := BaseFile(filepath.Base(src))
	log = filepath.Join(dir, dst)

	_, err = sys.StdCall("dotnet", cmd, "-f", src, "-q", "--csv", dir, "--csvf", dst+".csv")

	switch filepath.Ext(src) {
	case ".automaticDestinations-ms":
		log += "_AutomaticDestinations.csv"
	case ".customDestinations-ms":
		log += "_CustomDestinations.csv"
	}

	return
}

func Sbe(src, dir string) (log string, err error) {
	cmd, err := ez.Path("SBECmd.dll")

	if err != nil {
		return
	}

	if len(dir) == 0 {
		dir = filepath.Dir(src)
	}

	b := BaseFile(filepath.Base(src))

	dst := "out.csv"
	tmp := filepath.Join(dir, "tmp")
	log = filepath.Join(dir, fmt.Sprintf("%s_%s", b, dst))

	if err = os.MkdirAll(tmp, sys.MODE_DIR); err != nil {
		return
	}

	if err = Copy(tmp, src); err != nil {
		return
	}

	_, err = sys.StdCall("dotnet", cmd, "-d", tmp, "--csv", dir, "--csvf", dst)

	if err != nil {
		return
	}

	if err = os.Remove(filepath.Join(dir, "!SBECmd_Messages.txt")); err != nil {
		return
	}

	err = os.RemoveAll(tmp)

	return
}
