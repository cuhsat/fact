// FLog implementation details.
package flog

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cuhsat/fact/internal/fact"
	"github.com/cuhsat/fact/internal/flog"
	"github.com/cuhsat/fact/internal/sys"
	"github.com/cuhsat/fact/pkg/ecs"
)

const (
	Evtx = "evtx"
)

func LogEvent(src, dir string, pty bool) (logs []string, err error) {
	lines, err := flog.ImportEvent(src, dir)

	if err != nil {
		return
	}

	b := filepath.Base(src)

	f := strings.TrimSuffix(b, filepath.Ext(b))

	for i, line := range lines {
		dst := filepath.Join(dir, fmt.Sprintf("%s_%08d.json", f, i))

		e, err := ecs.MapEvent(line, src)

		if err != nil {
			sys.Error(err)
			continue
		}

		b, err := e.Bytes(pty)

		if err != nil {
			sys.Error(err)
			continue
		}

		if err = flog.ExportEvent(b, dst); err != nil {
			sys.Error(err)
			continue
		}

		l, err := filepath.Abs(dst)

		if err != nil {
			sys.Error(err)
			continue
		}

		if sys.Progress != nil {
			sys.Progress(l)
		}

		logs = append(logs, l)
	}

	return logs, nil
}

func StripHash(in []string) (out []string) {
	if len(in) == 0 {
		return in
	}

	i := strings.Index(in[0], fact.HashSep)

	if i == -1 {
		return in
	}

	for _, l := range in {
		out = append(out, l[i+2:])
	}

	return
}
