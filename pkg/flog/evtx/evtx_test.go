// Evtx implementation tests.
package evtx

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/cuhsat/fact/internal/fact/zip"
	"github.com/cuhsat/fact/internal/test"
)

func TestLog(t *testing.T) {
	cases := []struct {
		name, file, evtx, json string
	}{
		{
			name: "Test log for Windows",
			file: test.Testdata("windows", "event.zip"),
			evtx: "System.evtx",
			json: "System_00000000.json",
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

			err := Log(evt, tmp)

			if err != nil {
				t.Fatal(err)
			}

			f, err := files(tmp)

			if err != nil {
				t.Fatal(err)
			}

			if len(f) != 2 {
				t.Fatal("file count differs")
			}

			b, err := os.ReadFile(filepath.Join(tmp, tt.json))

			if err != nil {
				t.Fatal(err)
			}

			if !json.Valid(b) {
				t.Fatal("invalid json")
			}
		})
	}
}

func files(dir string) (f []string, err error) {
	de, err := os.ReadDir(dir)

	if err != nil {
		return
	}

	for _, e := range de {
		if !e.IsDir() {
			f = append(f, e.Name())
		}
	}

	return
}
