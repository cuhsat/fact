// File functions.
package flog

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuhsat/fact/internal/sys"
)

func BaseFile(name string) string {
	b := filepath.Base(name)

	return strings.TrimSuffix(b, filepath.Ext(b))
}

func Copy(dir, src string) (err error) {
	dst := filepath.Join(dir, filepath.Base(src))

	b, err := os.ReadFile(src)

	if err != nil {
		return
	}

	err = os.WriteFile(dst, b, sys.MODE_FILE)

	if os.IsNotExist(err) {
		err = nil
	}

	return
}

func ConsumeJson(name string) (lines []string, err error) {
	f, err := os.Open(name)

	if err != nil {
		return
	}

	fs := bufio.NewScanner(f)

	fs.Split(bufio.ScanLines)

	for fs.Scan() {
		lines = append(lines, fs.Text())
	}

	f.Close()

	err = os.Remove(name)

	return
}

func ConsumeCsv(name string) (lines []string, err error) {
	f, err := os.Open(name)

	if err != nil {
		return
	}

	rr, err := csv.NewReader(f).ReadAll()

	if len(rr) <= 1 {
		f.Close()
		return
	}

	m := map[string]string{}

	for _, r := range rr[1:] {
		for i, c := range r {
			m[rr[0][i]] = c
		}

		b, err := json.Marshal(m)

		if err != nil {
			sys.Error(err)
			continue
		}

		lines = append(lines, string(b))
	}

	f.Close()

	err = os.Remove(name)

	return
}
