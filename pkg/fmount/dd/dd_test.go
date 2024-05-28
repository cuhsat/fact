// DD implementation tests.
package dd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuhsat/fact/internal/fact/zip"
	"github.com/cuhsat/fact/internal/test"
)

func TestDD(t *testing.T) {
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

			if len(p) != 1 {
				t.Fatal("partition count differs")
			}

			sys := filepath.Join(mnt, "p1")

			if p[0] != sys {
				t.Fatal("mount point does not exist")
			}

			dir, _ := os.ReadDir(sys)

			if len(dir) == 0 {
				t.Fatal("mount point is empty")
			}

			err = Unmount(img)

			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
