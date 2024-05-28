// ECS event mapping functions.
package ecs

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"time"

	"github.com/cuhsat/fact/internal/fact"
	"github.com/cuhsat/fact/internal/flog"
)

type Event struct {
	Ecs     Ecs
	Agent   Agent
	File    File
	Base    Base
	Event   Evt
	Host    Host
	User    User
	Process Process
}

func (e *Event) Bytes() (b []byte, err error) {
	return json.Marshal(e)
}

func MapEvent(log, src string) (e *Event, err error) {
	m, err := flog.NewMap(log)

	if err != nil {
		return
	}

	timestamp := m.GetTime("Event/System/TimeCreated/@SystemTime")
	timezone, _ := timestamp.Zone()

	message := m.GetString("Event/EventData/Data/#text")

	hostname := m.GetString("Event/System/Computer")

	channel := m.GetString("Event/System/Channel")
	level := m.GetInt64("Event/System/Level")
	task := m.GetInt64("Event/System/Task")

	e = &Event{
		Ecs: Ecs{
			Version: Version,
		},
		Agent: Agent{
			Type:    fact.Product,
			Version: fact.Version,
		},
		File: File{
			Name:        filepath.Base(src),
			Directory:   filepath.Dir(src),
			Extension:   strings.Replace(filepath.Ext(src), ".", "", 1),
			DriveLetter: strings.Replace(filepath.VolumeName(src), ":", "", 1),
			Path:        src,
			Type:        "file",
		},
		Base: Base{
			Timestamp: timestamp,
			Message:   message,
			Tags:      "EventLog",
			Labels: map[string]interface{}{
				"Channel": channel,
				"Level":   level,
				"Task":    task,
			},
		},
		Event: Evt{
			Kind:     "event",
			Module:   "EventLog",
			Dataset:  "EventLog." + channel,
			Severity: level,
			ID:       m.GetString("Event/System/EventRecordID"),
			Code:     m.GetString("Event/System/EventID/#text"),
			Provider: m.GetString("Event/System/Provider/@Name"),
			Timezone: timezone,
			Created:  timestamp,
			Ingested: time.Now().UTC(),
			Original: log,
			Hash:     Fingerprint(log),
		},
		Host: Host{
			Hostname: hostname,
			Name:     hostname,
		},
		User: User{
			ID: m.GetString("Event/System/Security/@UserID"),
		},
		Process: Process{
			PID:      m.GetInt64("Event/System/Execution/@ProcessID"),
			ThreadID: m.GetInt64("Event/System/Execution/@ThreadID"),
		},
	}

	return
}
