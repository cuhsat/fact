// FLog implementation details.
package flog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuhsat/fact/internal/fact"
	"github.com/cuhsat/fact/internal/flog"
	"github.com/cuhsat/fact/internal/sys"
	"github.com/cuhsat/fact/pkg/ecs"
	"golang.org/x/sync/errgroup"
)

type fnlog func(string, string, bool) ([]string, error)

func Log(files []string, dir string, jp bool) error {
	g := new(errgroup.Group)

	for _, f := range files {
		var fn fnlog

		ext := strings.ToLower(filepath.Ext(f))

		if ext == "evtx" {
			fn = LogEvent
		} else if strings.HasSuffix(ext, "destinations-ms") {
			fn = LogJumpList
		} else {
			sys.Error("ignored", f)
			continue
		}

		g.Go(func() error {
			_, err := fn(f, dir, jp)
			return err
		})
	}

	return g.Wait()
}

func LogEvent(src, dir string, jp bool) (logs []string, err error) {
	log, err := flog.EvtxeCmd(src, dir)

	if err != nil {
		return
	}

	ll, err := flog.ConsumeJson(log)

	if err != nil {
		return
	}

	f := flog.BaseFile(src)

	for i, l := range ll {
		dst := filepath.Join(dir, fmt.Sprintf("%s_%08d.json", f, i))

		m, err := ecs.MapEvent(l, src)

		if err != nil {
			sys.Error(err)
			continue
		}

		log, err = write(m, dst, jp)

		if err != nil {
			sys.Error(err)
			continue
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func LogJumpList(src, dir string, jp bool) (logs []string, err error) {
	log, err := flog.JleCmd(src, dir)

	if err != nil {
		return
	}

	if _, err = os.Stat(log); os.IsNotExist(err) {
		return logs, nil
	}

	ll, err := flog.ConsumeCsv(log)

	if err != nil {
		return
	}

	f := flog.BaseFile(src)

	for i, l := range ll {
		dst := filepath.Join(dir, fmt.Sprintf("%s_%08d.json", f, i))

		m, err := ecs.MapJumpList(l, src)

		if err != nil {
			sys.Error(err)
			continue
		}

		log, err = write(m, dst, jp)

		if err != nil {
			sys.Error(err)
			continue
		}

		logs = append(logs, log)
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

func write(a any, dst string, jp bool) (log string, err error) {
	var b []byte

	if jp {
		b, err = json.MarshalIndent(a, "", "  ")
	} else {
		b, err = json.Marshal(a)
	}

	if err != nil {
		return
	}

	f, err := os.Create(dst)

	if err != nil {
		return
	}

	defer f.Close()

	_, err = f.Write(b)

	if err != nil {
		return
	}

	log, err = filepath.Abs(dst)

	if err == nil && sys.Progress != nil {
		sys.Progress(log)
	}

	return
}
