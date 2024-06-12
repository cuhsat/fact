// Eric Zimmermann tools.
package flog

import (
	"path/filepath"

	"github.com/cuhsat/fact/internal/fact/ez"
	"github.com/cuhsat/fact/internal/sys"
)

func EvtxeCmd(src, dir string) (log string, err error) {
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

func JleCmd(src, dir string) (log string, err error) {
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
