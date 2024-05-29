// Zip archive implementation tests.
package zip

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/cuhsat/fact/internal/test"
)

func TestIndex(t *testing.T) {
	t.Run("Test index", func(t *testing.T) {
		idx, err := Index(test.Testdata("windows", "image.zip"))

		if err != nil {
			t.Fatal(err)
		}

		if !slices.Contains(idx, "Users/Test/NTUSER.DAT") {
			t.Fatal("file not found")
		}
	})
}

func TestUnzip(t *testing.T) {
	t.Run("Test unzip", func(t *testing.T) {
		tmp, _ := os.MkdirTemp(os.TempDir(), "zip")

		err := Unzip(test.Testdata("windows", "image.zip"), tmp)

		if err != nil {
			t.Fatal(err)
		}

		_, err = os.Stat(filepath.Join(tmp, "Windows"))

		if os.IsNotExist(err) {
			t.Fatal("folder not found")
		}
	})
}
