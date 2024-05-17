// Test functions.
package test

import (
	"path"
	"runtime"
)

func Testdata(name string) string {
	_, c, _, ok := runtime.Caller(0)

	if !ok {
		return "error"
	}

	return path.Join(path.Dir(c), "..", "testdata", name)
}
