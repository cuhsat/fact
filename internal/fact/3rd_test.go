// Fact 3rd party tests.
package fact

import (
	"os"
	"testing"
)

func TestTools(t *testing.T) {
	os.Setenv("EZTOOLS", "../../bin")

	cases := []struct {
		name, tool string
	}{
		{
			name: "Test for EvtxECmd",
			tool: "EvtxECmd.dll",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			p, err := EzTools(tt.tool)

			if err != nil {
				t.Fatal(err)
			}

			if len(p) == 0 {
				t.Fatal(tt.tool + " not found")
			}

			if _, err := os.Stat(p); os.IsNotExist(err) {
				t.Fatal(tt.tool + " not found")
			}
		})
	}
}
