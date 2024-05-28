// ECS specification.
package ecs

import (
	"crypto/sha1"
	"encoding/hex"
	"time"
)

const (
	Version = "8.11"
)

type Ecs struct {
	Version string `ecs:"version"`
}

type Agent struct {
	Type    string `ecs:"type"`
	Version string `ecs:"version"`
}

type File struct {
	Name        string `ecs:"name"`
	Directory   string `ecs:"directory"`
	Extension   string `ecs:"extension"`
	DriveLetter string `ecs:"drive_letter"`
	Path        string `ecs:"path"`
	Type        string `ecs:"type"`
}

type Base struct {
	Timestamp time.Time              `ecs:"@timestamp"`
	Message   string                 `ecs:"message"`
	Tags      string                 `ecs:"tags"`
	Labels    map[string]interface{} `ecs:"labels"`
}

type Evt struct {
	Kind     string    `ecs:"kind"`
	Module   string    `ecs:"module"`
	Dataset  string    `ecs:"dataset"`
	Severity int64     `ecs:"severity"`
	ID       string    `ecs:"id"`
	Code     string    `ecs:"code"`
	Provider string    `ecs:"provider"`
	Timezone string    `ecs:"timezone"`
	Created  time.Time `ecs:"created"`
	Ingested time.Time `ecs:"ingested"`
	Original string    `ecs:"original"`
	Hash     string    `ecs:"hash"`
}

type Host struct {
	Hostname string `ecs:"hostname"`
	Name     string `ecs:"name"`
}

type User struct {
	ID string `ecs:"id"`
}

type Process struct {
	PID      int64 `ecs:"pid"`
	ThreadID int64 `ecs:"thread.id"`
}

func Fingerprint(s string) string {
	h := sha1.New()
	h.Write([]byte(s))

	return hex.EncodeToString(h.Sum(nil))
}
