// Evtx implementation details.
package evtx

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuhsat/fact/internal/fact/ez"
	"github.com/cuhsat/fact/internal/flog"
	"github.com/cuhsat/fact/internal/sys"
	"github.com/cuhsat/fact/pkg/ecs"
)

const (
	Evtx = "evtx"
)

func Log(src string, dir string) (err error) {
	lines, err := _import(src, dir)

	if err != nil {
		return
	}

	b := filepath.Base(src)

	f := strings.TrimSuffix(b, filepath.Ext(b))

	for i, line := range lines {
		dst := filepath.Join(dir, fmt.Sprintf("%s_%08d.json", f, i))

		if err = _export(src, dst, line); err != nil {
			sys.Error(err)
			continue
		}
	}

	return nil
}

func _import(src, dir string) (lines []string, err error) {
	log, err := evtxecmd(src, src+".json", dir)

	if err != nil {
		return
	}

	lines, err = flog.ReadLines(log)

	if err != nil {
		return
	}

	err = os.Remove(log)

	return
}

func _export(src, dst, line string) (err error) {
	e, err := ecs.MapEvent(line, src)

	if err != nil {
		return
	}

	b, err := e.Bytes()

	if err != nil {
		return
	}

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

	_, err = sys.StdCall("dotnet", asm, "-f", src, "--fj", "--json", dir, "--jsonf", dst)

	log = filepath.Join(dir, dst)

	return
}
