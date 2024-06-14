// FLog implementation details.
package flog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/cuhsat/fact/internal/fact"
	"github.com/cuhsat/fact/internal/flog"
	"github.com/cuhsat/fact/internal/sys"
	"github.com/cuhsat/fact/pkg/ecs"
	"golang.org/x/sync/errgroup"
)

type fnlog func(string, string, bool) ([]string, error)

func Log(files []string, dir string, jp bool) error {
	usrHives := []string{
		"ntuser.dat",
		"usrclass.dat",
	}

	usrHistory := []string{
		"history",
		"places.sqlite",
	}

	g := new(errgroup.Group)

	for _, f := range files {
		var fn fnlog

		name := strings.ToLower(filepath.Base(f))
		ext := strings.ToLower(filepath.Ext(f))

		if ext == ".evtx" {
			fn = LogEvent
		} else if strings.HasSuffix(ext, "destinations-ms") {
			fn = LogJumpList
		} else if slices.Contains(usrHives, name) {
			fn = LogShellBag
		} else if slices.Contains(usrHistory, name) {
			fn = LogHistory
		} else {
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
	log, err := flog.Evtxe(src, dir)

	if err != nil {
		return
	}

	ll, err := flog.ConsumeJson(log)

	if err != nil {
		return
	}

	for _, l := range ll {
		dst := filepath.Join(dir, fmt.Sprintf("%s.json", ecs.Hash(l)))

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
	log, err := flog.Jle(src, dir)

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

	for _, l := range ll {
		dst := filepath.Join(dir, fmt.Sprintf("%s.json", ecs.Hash(l)))

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

func LogShellBag(src, dir string, jp bool) (logs []string, err error) {
	log, err := flog.Sbe(src, dir)

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

	for _, l := range ll {
		dst := filepath.Join(dir, fmt.Sprintf("%s.json", ecs.Hash(l)))

		m, err := ecs.MapShellBag(l, src)

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

func LogHistory(src, dir string, jp bool) (logs []string, err error) {
	ll, err := flog.History(src)

	if err != nil {
		return
	}

	if err = os.MkdirAll(dir, sys.MODE_DIR); err != nil {
		return
	}

	for _, l := range ll {
		dst := filepath.Join(dir, fmt.Sprintf("%s.json", ecs.Hash(fmt.Sprint(l))))

		m, err := ecs.MapHistory(&l, src)

		if err != nil {
			sys.Error(err)
			continue
		}

		log, err := write(m, dst, jp)

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
