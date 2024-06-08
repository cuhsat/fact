// FLog implementation tests.
package flog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/cuhsat/fact/internal/fact"
	"github.com/cuhsat/fact/internal/fact/zip"
	"github.com/cuhsat/fact/internal/test"
)

func TestMain(m *testing.M) {
	os.Setenv("EZTOOLS", "../../../bin")
	os.Exit(m.Run())
}

func TestLog(t *testing.T) {
	cases := []struct {
		name, file, evtx string
	}{
		{
			name: "Test log for Windows",
			file: test.Testdata("windows", "evtx.zip"),
			evtx: "System.evtx",
		},
	}

	for _, tt := range cases {
		tmp, _ := os.MkdirTemp(os.TempDir(), "log")

		err := zip.Unzip(tt.file, tmp)

		if err != nil {
			t.Fatal(err)
		}

		t.Run(tt.name, func(t *testing.T) {
			evt := filepath.Join(tmp, tt.evtx)

			l, err := Evtx(evt, tmp, true)

			if err != nil {
				t.Fatal(err)
			}

			if len(l) != 1 {
				t.Fatal("file count differs")
			}

			b, err := os.ReadFile(l[0])

			if err != nil {
				t.Fatal(err)
			}

			if !json.Valid(b) {
				t.Fatal("invalid json")
			}
		})
	}
}

func TestStripHash(t *testing.T) {
	cases := []struct {
		name, file, hash string
	}{
		{
			name: "Test StripHash",
			file: "test",
			hash: "68ac906495480a3404beee4874ed853a037a7a8f",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			f := StripHash([]string{
				tt.hash + fact.HashSep + tt.file,
			})

			if len(f) != 1 {
				t.Fatal("file count wrong")
			}

			if f[0] != tt.file {
				t.Fatal("hash not stripped")
			}
		})
	}
}
