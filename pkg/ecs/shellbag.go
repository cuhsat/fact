// ECS shellbag mapping functions.
package ecs

import (
	"github.com/cuhsat/fact/internal/flog"
)

func MapShellBag(s, src string) (log *Log, err error) {
	m, err := flog.NewMap(s)

	if err != nil {
		return
	}

	log = NewLog(s, src, &Base{
		Timestamp: m.GetTime("LastInteracted", "LastWriteTime"),
		Message:   m.GetString("AbsolutePath"),
		Tags:      "ShellBag",
	})

	log.Registry = &Registry{
		Hive: "HKU",
	}

	return
}
