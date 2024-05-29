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
			file: test.Testdata("windows", "event.zip"),
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

			l, err := Log(evt, tmp, true)

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
