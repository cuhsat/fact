// ECS event mapping functions.
package ecs

import (
	"time"

	"github.com/cuhsat/fact/internal/flog"
)

func MapEvent(s, src string) (log *Log, err error) {
	m, err := flog.NewMap(s)

	if err != nil {
		return
	}

	channel := m.GetString("Event/System/Channel")

	timestamp := m.GetTime("Event/System/TimeCreated/@SystemTime")
	timezone, _ := timestamp.Zone()
	message := m.GetString("Event/EventData/Data/#text")

	return NewLog(
		src,
		Base{
			Timestamp: timestamp,
			Message:   message,
			Tags:      "EventLog",
			Labels: map[string]interface{}{
				"Channel": channel,
				"Level":   m.GetInt64("Event/System/Level"),
				"Task":    m.GetInt64("Event/System/Task"),
			},
		},
		Evt{
			Kind:     "event",
			Module:   "EventLog",
			Dataset:  "EventLog." + channel,
			Severity: m.GetInt64("Event/System/Level"),
			ID:       m.GetString("Event/System/EventRecordID"),
			Code:     m.GetString("Event/System/EventID/#text"),
			Provider: m.GetString("Event/System/Provider/@Name"),
			Timezone: timezone,
			Created:  timestamp,
			Ingested: time.Now().UTC(),
			Original: s,
			Hash:     GetHash(s),
		},
		Host{
			Hostname: m.GetString("Event/System/Computer"),
			Name:     m.GetString("Event/System/Computer"),
		},
		User{
			ID: m.GetString("Event/System/Security/@UserID"),
		},
		Process{
			PID:      m.GetInt64("Event/System/Execution/@ProcessID"),
			ThreadID: m.GetInt64("Event/System/Execution/@ThreadID"),
		},
	), nil
}
