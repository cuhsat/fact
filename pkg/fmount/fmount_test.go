// FMount implementation tests.
package fmount

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuhsat/fact/internal/fact/zip"
	"github.com/cuhsat/fact/internal/fmount"
	"github.com/cuhsat/fact/internal/test"
)

func TestMount(t *testing.T) {
	cases := []struct {
		name, file, path string
	}{
		{
			name: "Test mount for Windows (dd)",
			file: test.Testdata("windows", "disk.dd.zip"),
			path: "dd",
		},
		{
			name: "Test mount for Windows (vmdk)",
			file: test.Testdata("windows", "disk.vmdk.zip"),
			path: "dd",
		},
	}

	for _, tt := range cases {
		tmp, _ := os.MkdirTemp(os.TempDir(), tt.path)
		mnt, _ := os.MkdirTemp(os.TempDir(), tt.path+"-mnt")

		err := zip.Unzip(tt.file, tmp)

		if err != nil {
			t.Fatal(err)
		}

		t.Run(tt.name, func(t *testing.T) {
			img := filepath.Join(tmp, fmount.BaseFile(tt.file))

			p, err := Mount(img, mnt, "", true)

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

			if _, err = os.Stat(mnt); !os.IsNotExist(err) {
				t.Fatal("mount point not removed")
			}
		})
	}
}
