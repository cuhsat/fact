// ECS jumplist mapping functions.
package ecs

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/cuhsat/fact/internal/flog"
)

func MapJumpList(s, src string) (log *Log, err error) {
	m, err := flog.NewMap(s)

	if err != nil {
		return
	}

	exec := m.GetString("LocalPath", "Path")
	arg := m.GetString("Arguments")

	var args []string

	if len(arg) > 0 {
		args = strings.Split(arg, " ")
	}

	timestamp := m.GetTime("LastModified", "TargetAccessed")
	timezone, _ := timestamp.Zone()
	message := strings.TrimSpace(exec + " " + m.GetString("Arguments"))

	return NewLog(
		src,
		Base{
			Timestamp: timestamp,
			Message:   message,
			Tags:      "JumpList",
			Labels: map[string]interface{}{
				"Destination": target(s),
			},
		},
		Evt{
			Timezone: timezone,
			Created:  timestamp,
			Ingested: time.Now().UTC(),
			Original: s,
			Hash:     GetHash(s),
		},
		Host{
			Hostname: m.GetString("Hostname", "MachineID"),
			Name:     m.GetString("Hostname", "MachineID"),
			MAC:      m.GetString("MacAddress", "MachineMACAddress"),
		},
		User{},
		Process{
			EntityID:         m.GetString("AppId"),
			Start:            timestamp,
			Name:             filepath.Base(exec),
			Title:            m.GetString("AppIdDescription"),
			Executable:       exec,
			Args:             args,
			ArgsCount:        int64(len(args)),
			CommandLine:      message,
			WorkingDirectory: m.GetString("WorkingDirectory"),
		},
	), nil
}

func target(log string) string {
	if strings.Contains(log, "DestListVersion") {
		return "automatic"
	} else {
		return "custom"
	}
}
