// Fact 3rd party functions.
package fact

import (
	"os"
	"path/filepath"
)

func EzTools(asm string) (p string, err error) {
	env := os.ExpandEnv("$EZTOOLS")

	if len(env) > 0 {
		p = filepath.Join(env, asm)
		return
	}

	wd, err := os.Getwd()

	p = filepath.Join(wd, asm)

	return
}
