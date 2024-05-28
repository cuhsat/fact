// Map functions.
package flog

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/cuhsat/fact/internal/sys"
)

const (
	KeySep = "/"
	ValSep = "\n"
)

var (
	bom = []byte{0xEF, 0xBB, 0xBF}
)

type Map struct {
	o object
}

type object any

func NewMap(s string) (m *Map, err error) {
	m = &Map{}

	return m, json.Unmarshal(bytes.TrimPrefix([]byte(s), bom), &m.o)
}

func (m *Map) GetString(key string) (value string) {
	return rget(m.o, strings.Split(key, KeySep))
}

func (m *Map) GetInt64(key string) (value int64) {
	value, err := strconv.ParseInt(m.GetString(key), 10, 64)

	if err != nil {
		return -1 // default
	}

	return
}

func (m *Map) GetTime(key string) (value time.Time) {
	const layout = "2006-01-02 15:04:05.9999999"

	value, err := time.Parse(layout, m.GetString(key))

	if err != nil {
		return time.UnixMicro(0) // default
	}

	return
}

func rget(o object, k []string) (s string) {
	switch v := o.(type) {
	case map[string]any:
		o, ok := v[k[0]]

		if !ok {
			return
		}

		if len(k) > 1 {
			return rget(o, k[1:])
		}

		s, ok = o.(string)

		if !ok {
			sys.Error("map: key too short")
			return
		}

		return

	case []any:
		var l []string

		for _, o = range v {
			l = append(l, rget(o, k))
		}

		return strings.Join(l, ValSep)

	case string: // abort for cut short keys
		return v

	case nil:
		return
	}

	sys.Error("map: type not supported")
	return
}
