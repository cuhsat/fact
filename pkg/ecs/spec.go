// ECS specification.
package ecs

import (
	"crypto/sha1"
	"encoding/hex"
	"path/filepath"
	"strings"
	"time"

	"github.com/cuhsat/fact/internal/fact"
)

const (
	Version = "8.11"
)

type Log struct {
	Ecs     Ecs     `json:"ecs"`
	Agent   Agent   `json:"agent"`
	Base    Base    `json:"base"`
	File    File    `json:"file,omitempty"`
	Event   Evt     `json:"event,omitempty"`
	Host    Host    `json:"host,omitempty"`
	User    User    `json:"user,omitempty"`
	Process Process `json:"process,omitempty"`
}

type Ecs struct {
	Version string `json:"version"`
}

type Agent struct {
	Type    string `json:"type"`
	Version string `json:"version"`
}

type Base struct {
	Timestamp time.Time              `json:"@timestamp"`
	Message   string                 `json:"message"`
	Tags      string                 `json:"tags,omitempty"`
	Labels    map[string]interface{} `json:"labels,omitempty"`
}

type File struct {
	Name        string `json:"name,omitempty"`
	Directory   string `json:"directory,omitempty"`
	Extension   string `json:"extension,omitempty"`
	DriveLetter string `json:"drive_letter,omitempty"`
	Path        string `json:"path,omitempty"`
	Type        string `json:"type,omitempty"`
}

type Evt struct {
	Kind     string    `json:"kind,omitempty"`
	Module   string    `json:"module,omitempty"`
	Dataset  string    `json:"dataset,omitempty"`
	Severity int64     `json:"severity,omitempty"`
	ID       string    `json:"id,omitempty"`
	Code     string    `json:"code,omitempty"`
	Provider string    `json:"provider,omitempty"`
	Timezone string    `json:"timezone,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	Ingested time.Time `json:"ingested,omitempty"`
	Original string    `json:"original,omitempty"`
	Hash     string    `json:"hash,omitempty"`
}

type Host struct {
	Hostname string `json:"hostname,omitempty"`
	Name     string `json:"name,omitempty"`
	MAC      string `json:"mac,omitempty"`
}

type User struct {
	ID string `json:"id,omitempty"`
}

type Process struct {
	PID              int64     `json:"pid,omitempty"`
	ThreadID         int64     `json:"thread.id,omitempty"`
	EntityID         string    `json:"entity_id,omitempty"`
	Start            time.Time `json:"start,omitempty"`
	Name             string    `json:"name,omitempty"`
	Title            string    `json:"title,omitempty"`
	Args             []string  `json:"args,omitempty"`
	ArgsCount        int64     `json:"args_count,omitempty"`
	Executable       string    `json:"executable,omitempty"`
	CommandLine      string    `json:"command_line,omitempty"`
	WorkingDirectory string    `json:"working_directory,omitempty"`
}

func NewLog(src string, base Base, event Evt, host Host, user User, process Process) *Log {
	return &Log{
		Ecs: Ecs{
			Version: Version,
		},
		Agent: Agent{
			Type:    fact.Product,
			Version: fact.Version,
		},
		Base: base,
		File: File{
			Name:        filepath.Base(src),
			Directory:   r(filepath.Abs(filepath.Dir(src))),
			Extension:   strings.Replace(filepath.Ext(src), ".", "", 1),
			DriveLetter: strings.Replace(filepath.VolumeName(src), ":", "", 1),
			Path:        r(filepath.Abs(src)),
			Type:        "file",
		},
		Event:   event,
		Host:    host,
		User:    user,
		Process: process,
	}
}

func GetHash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))

	return hex.EncodeToString(h.Sum(nil))
}

func r(a string, _ error) string {
	return a
}
