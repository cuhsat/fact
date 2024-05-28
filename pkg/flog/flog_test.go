// FLog implementation tests.
package flog

import (
	"testing"

	"github.com/cuhsat/fact/internal/fact"
)

func TestStripHash(t *testing.T) {
	cases := []struct {
		name, file, hash string
	}{
		{
			name: "Test StripHash",
			file: "test",
			hash: "68ac906495480a3404beee4874ed853a037a7a8f",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			f := StripHash([]string{
				tt.hash + fact.HashSep + tt.file,
			})

			if len(f) != 1 {
				t.Fatal("file count wrong")
			}

			if f[0] != tt.file {
				t.Fatal("hash not stripped")
			}
		})
	}
}
