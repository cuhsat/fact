// FLog implementation details.
package flog

import (
	"path/filepath"
	"strings"

	"github.com/cuhsat/fact/internal/fact"
	"github.com/cuhsat/fact/internal/sys"
	"github.com/cuhsat/fact/pkg/flog/evt"
)

func Evt(files, args []string) (err error) {
	for _, f := range files {
		if filepath.Ext(f) == evt.Evt {
			args = append(args, f)
		}
	}

	_, err = sys.StdCall("flog.evt", args...)

	return
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