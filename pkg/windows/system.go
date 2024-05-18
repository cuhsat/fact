// Windows system artifact enumeration functions.
package windows

import (
	"os"
	"path/filepath"

	"github.com/cuhsat/fact/internal/sys"
)

func EnumSystem(sysroot string, out chan<- string) {
	if len(sysroot) == 0 {
		sysroot = os.ExpandEnv("$SYSTEMDRIVE")
	}

	root := filepath.Join(sysroot, "Windows")

	if _, err := os.Stat(root); err != nil {
		sys.Error(err)
	}

	for _, artifact := range [...]string{
		"[Ss]ystem32/[Cc]onfig/[Cc][Oo][Mm][Pp][Oo][Nn][Ee][Nn][Tt][Ss]",
		"[Ss]ystem32/[Cc]onfig/[Dd][Ee][Ff][Aa][Uu][Ll][Tt]",
		"[Ss]ystem32/[Cc]onfig/[Ss][Aa][Mm]",
		"[Ss]ystem32/[Cc]onfig/[Ss][Ee][Cc][Uu][Rr][Ii][Tt][Yy]",
		"[Ss]ystem32/[Cc]onfig/[Ss][Oo][Ff][Tt][Ww][Aa][Rr][Ee]",
		"[Ss]ystem32/[Cc]onfig/[Ss][Yy][Ss][Tt][Ee][Mm]",
		"[Ss]ystem32/[Ww]inevt/[Ll]ogs/*.evt*",
		"[Pp]refetch/*.pf",
		"[Aa]m[Cc]ompat/[Pp]rograms/[Aa]m[Cc]ache.hve",
	} {
		files, err := filepath.Glob(filepath.Join(root, artifact))

		if err != nil {
			sys.Error(err)
			continue
		}

		for _, file := range files {
			out <- filepath.ToSlash(file)
		}
	}
}
