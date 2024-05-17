// System implementation tests.
package sys

import (
	"strings"
	"testing"
)

func TestCall(t *testing.T) {
	t.Run("Test layered call", func(t *testing.T) {
		stdout := new(strings.Builder)
		stderr := new(strings.Builder)

		code := call(stdout, stderr, "go", "version")

		if !strings.Contains(stdout.String(), "go version go") {
			t.Fatal("Stdout", stdout.String())
		}

		if len(stderr.String()) > 0 {
			t.Fatal("Stderr", stderr.String())
		}

		if code != 0 {
			t.Fatal("Exit code", code)
		}
	})
}
