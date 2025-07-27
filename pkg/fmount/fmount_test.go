// FMount implementation tests.
package fmount

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cuhsat/fact/internal/fact/zip"
	"github.com/cuhsat/fact/internal/fmount"
	"github.com/cuhsat/fact/internal/sys"
	"github.com/cuhsat/fact/internal/test"
)

func TestMain(m *testing.M) {
	sys.Progress = nil

	if _, ci := os.LookupEnv("GITHUB_ACTIONS"); !ci {
		os.Exit(m.Run())
	}
}

func TestMount(t *testing.T) {
	cases := []struct {
		name, file, path string
	}{
		{
			name: "Test mount with disk image",
			file: test.Testdata("windows", "disk.zip"),
			path: "disk",
		},
	}

	for _, tt := range cases {
		tmp, err := os.MkdirTemp(os.TempDir(), tt.path)

		if err != nil {
			t.Fatal(err)
		}

		mnt, err := os.MkdirTemp(os.TempDir(), tt.path+"-mnt")

		if err != nil {
			t.Fatal(err)
		}

		err = zip.Unzip(tt.file, tmp)

		if err != nil {
			t.Fatal(err)
		}

		t.Run(tt.name, func(t *testing.T) {
			img := filepath.Join(tmp, fmount.BaseFile(tt.file))

			p, err := Mount(img, mnt, "", true, []string{})

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
