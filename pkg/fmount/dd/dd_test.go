// DD implementation tests.
package dd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuhsat/fact/internal/test"
	"github.com/cuhsat/fact/internal/zip"
)

func TestDD(t *testing.T) {
	if len(os.Getenv("CI")) > 0 {
		t.Skip("skip test")
	}

	cases := []struct {
		name, file string
	}{
		{
			name: "Test mount for Windows",
			file: test.Testdata("windows.dd.zip"),
		},
	}

	for _, tt := range cases {
		tmp, _ := os.MkdirTemp(os.TempDir(), "dd")
		mnt, _ := os.MkdirTemp(os.TempDir(), "dd-mnt")

		err := zip.Unzip(tt.file, tmp)

		if err != nil {
			t.Fatal(err)
		}

		t.Run(tt.name, func(t *testing.T) {
			img := filepath.Join(tmp, baseFile(tt.file))

			p, err := Mount(img, mnt, true)

			if err != nil {
				t.Fatal(err)
			}

			if p != filepath.Join(mnt, "p1") {
				t.Fatal("sysroot unexpected", p)
			}

			err = Unmount(img)

			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
