// Test functions.
package test

import (
	"path/filepath"
	"runtime"
)

func Testdata(name string) string {
	_, c, _, ok := runtime.Caller(0)

	if !ok {
		return "error"
	}

	return filepath.Join(filepath.Dir(c), "..", "testdata", name)
}
