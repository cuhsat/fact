// FMount implementation tests.
package fmount

import (
	"testing"

	"github.com/cuhsat/fact/internal/test"
	"github.com/cuhsat/fact/pkg/fmount/dd"
)

func TestDetectType(t *testing.T) {
	cases := []struct {
		name, img, val string
	}{
		{
			name: "Test with type",
			img:  "win",
			val:  "dd",
		},
		{
			name: "Test with extension",
			img:  "win.dd",
		},
		{
			name: "Test without anything",
			img:  test.Testdata("mbr"),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			v, err := DetectType(tt.img, tt.val)

			if err != nil {
				t.Fatal(err)
			}

			if v != dd.DD && v != dd.RAW {
				t.Fatal("Type mismatch", v)
			}
		})
	}
}
